package bot

import (
	"fmt"
	"sync"
	"time"

	"github.com/feedme/order-controller/internal/order"
	"github.com/feedme/order-controller/internal/queue"
)

type BotStatus int

const (
	Idle BotStatus = iota
	BotProcessing
)

type Bot struct {
	ID              int
	Status          BotStatus
	CurrentOrder    *order.Order
	ProcessDuration time.Duration
	stopChan        chan struct{}
	doneChan        chan struct{}
	queue           *queue.OrderQueue
	orderMgr        *order.OrderManager
	notifyNew       chan struct{}
	mu              sync.Mutex
}

func NewBot(id int, q *queue.OrderQueue, om *order.OrderManager, notifyNew chan struct{}) *Bot {
	return &Bot{
		ID:              id,
		Status:          Idle,
		ProcessDuration: 10 * time.Second,
		queue:           q,
		orderMgr:        om,
		notifyNew:       notifyNew,
		stopChan:        make(chan struct{}),
		doneChan:        make(chan struct{}),
	}
}

func (b *Bot) Run() {
	defer close(b.doneChan)
	for {
		o := b.queue.Dequeue()
		if o != nil {
			b.mu.Lock()
			b.Status = BotProcessing
			b.CurrentOrder = o
			b.mu.Unlock()
			b.orderMgr.UpdateStatus(o.ID, order.Processing)
			select {
			case <-time.After(b.ProcessDuration):
				b.orderMgr.UpdateStatus(o.ID, order.Complete)
				b.mu.Lock()
				b.CurrentOrder = nil
				b.Status = Idle
				b.mu.Unlock()
			case <-b.stopChan:
				b.queue.RequeueAt(b.CurrentOrder, 0)
				b.orderMgr.UpdateStatus(b.CurrentOrder.ID, order.Pending)
				b.mu.Lock()
				b.CurrentOrder = nil
				b.Status = Idle
				b.mu.Unlock()
				return
			}
		} else {
			select {
			case <-b.stopChan:
				return
			case <-b.notifyNew:
			}
		}
	}
}

func (b *Bot) Stop() {
	close(b.stopChan)
	<-b.doneChan
}

func (b *Bot) GetStatus() (BotStatus, *order.Order) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Status, b.CurrentOrder
}

func (b *Bot) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.Status == BotProcessing && b.CurrentOrder != nil {
		return fmt.Sprintf("Bot#%d (PROCESSING: Order#%d)", b.ID, b.CurrentOrder.ID)
	}
	return fmt.Sprintf("Bot#%d (IDLE)", b.ID)
}

type BotManager struct {
	bots   []*Bot
	nextID int
	mu     sync.Mutex
}

func NewBotManager() *BotManager {
	return &BotManager{
		bots: make([]*Bot, 0),
	}
}

func (bm *BotManager) AddBot(q *queue.OrderQueue, om *order.OrderManager, notifyNew chan struct{}) *Bot {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.nextID++
	b := NewBot(bm.nextID, q, om, notifyNew)
	bm.bots = append(bm.bots, b)
	go b.Run()
	return b
}

func (bm *BotManager) RemoveBot() *Bot {
	bm.mu.Lock()
	if len(bm.bots) == 0 {
		bm.mu.Unlock()
		return nil
	}
	b := bm.bots[len(bm.bots)-1]
	bm.bots = bm.bots[:len(bm.bots)-1]
	bm.mu.Unlock()

	b.Stop()

	return b
}

func (bm *BotManager) GetBots() []*Bot {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	result := make([]*Bot, len(bm.bots))
	copy(result, bm.bots)
	return result
}

func (bm *BotManager) BotCount() int {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	return len(bm.bots)
}

func (bm *BotManager) NotifyAll() {
	bm.mu.Lock()
	bots := make([]*Bot, len(bm.bots))
	copy(bots, bm.bots)
	bm.mu.Unlock()

	for _, b := range bots {
		b.mu.Lock()
		if b.Status == Idle {
			select {
			case b.notifyNew <- struct{}{}:
			default:
			}
		}
		b.mu.Unlock()
	}
}
