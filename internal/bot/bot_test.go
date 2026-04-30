package bot

import (
	"testing"
	"time"

	"github.com/feedme/order-controller/internal/order"
	"github.com/feedme/order-controller/internal/queue"
)

func TestBot_IdleWhenNoOrder(t *testing.T) {
	t.Parallel()
	om := order.NewOrderManager()
	q := queue.NewOrderQueue()
	notifyNew := make(chan struct{}, 100)

	bm := NewBotManager()
	b := bm.AddBot(q, om, notifyNew)

	time.Sleep(50 * time.Millisecond)

	status, curOrder := b.GetStatus()
	if status != Idle {
		t.Errorf("expected bot to be Idle, got %d", status)
	}
	if curOrder != nil {
		t.Error("expected no current order for idle bot")
	}

	bm.RemoveBot()
}

func TestBot_ProcessOrder(t *testing.T) {
	om := order.NewOrderManager()
	q := queue.NewOrderQueue()
	notifyNew := make(chan struct{}, 100)

	bm := NewBotManager()
	b := bm.AddBot(q, om, notifyNew)
	b.ProcessDuration = 100 * time.Millisecond

	o := om.CreateOrder(order.Normal, "test")
	q.Enqueue(o)

	select {
	case notifyNew <- struct{}{}:
	default:
	}

	time.Sleep(50 * time.Millisecond)

	status, curOrder := b.GetStatus()
	if status != BotProcessing {
		t.Errorf("expected bot to be BotProcessing, got %d", status)
	}
	if curOrder == nil || curOrder.ID != o.ID {
		t.Error("expected bot to be processing the order")
	}

	if status, ok := om.GetOrderStatus(o.ID); !ok || status != order.Processing {
		t.Errorf("expected order status Processing, got %d, ok=%v", status, ok)
	}

	time.Sleep(150 * time.Millisecond)

	status, curOrder = b.GetStatus()
	if status != Idle {
		t.Errorf("expected bot to be Idle after processing, got %d", status)
	}
	if curOrder != nil {
		t.Error("expected no current order after processing")
	}

	if status, ok := om.GetOrderStatus(o.ID); !ok || status != order.Complete {
		t.Errorf("expected order status Complete, got %d, ok=%v", status, ok)
	}

	bm.RemoveBot()
}

func TestBot_StopWhileProcessing(t *testing.T) {
	om := order.NewOrderManager()
	q := queue.NewOrderQueue()
	notifyNew := make(chan struct{}, 100)

	bm := NewBotManager()
	b := bm.AddBot(q, om, notifyNew)

	o := om.CreateOrder(order.Normal, "test")
	q.Enqueue(o)

	select {
	case notifyNew <- struct{}{}:
	default:
	}

	time.Sleep(100 * time.Millisecond)

	status, _ := b.GetStatus()
	if status != BotProcessing {
		t.Skip("bot not processing yet, skipping stop test")
	}

	bm.RemoveBot()

	time.Sleep(100 * time.Millisecond)

	if status, ok := om.GetOrderStatus(o.ID); !ok || status != order.Pending {
		t.Errorf("expected order status Pending after bot stop, got %d, ok=%v", status, ok)
	}

	pending := q.PendingOrders()
	found := false
	for _, p := range pending {
		if p.ID == o.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected order to be back in queue after bot stop")
	}
}

func TestBotManager_AddBot(t *testing.T) {
	t.Parallel()
	om := order.NewOrderManager()
	q := queue.NewOrderQueue()
	notifyNew := make(chan struct{}, 100)

	bm := NewBotManager()

	b1 := bm.AddBot(q, om, notifyNew)
	if b1 == nil {
		t.Fatal("expected bot to be created")
	}
	if b1.ID != 1 {
		t.Errorf("expected bot ID 1, got %d", b1.ID)
	}
	if bm.BotCount() != 1 {
		t.Errorf("expected 1 bot, got %d", bm.BotCount())
	}

	b2 := bm.AddBot(q, om, notifyNew)
	if b2.ID != 2 {
		t.Errorf("expected bot ID 2, got %d", b2.ID)
	}
	if bm.BotCount() != 2 {
		t.Errorf("expected 2 bots, got %d", bm.BotCount())
	}

	bm.RemoveBot()
	bm.RemoveBot()
}

func TestBotManager_RemoveBot(t *testing.T) {
	t.Parallel()
	om := order.NewOrderManager()
	q := queue.NewOrderQueue()
	notifyNew := make(chan struct{}, 100)

	bm := NewBotManager()
	bm.AddBot(q, om, notifyNew)
	b2 := bm.AddBot(q, om, notifyNew)

	removed := bm.RemoveBot()
	if removed == nil {
		t.Fatal("expected bot to be removed")
	}
	if removed.ID != b2.ID {
		t.Errorf("expected newest bot (ID %d) to be removed, got ID %d", b2.ID, removed.ID)
	}
	if bm.BotCount() != 1 {
		t.Errorf("expected 1 bot after removal, got %d", bm.BotCount())
	}

	bm.RemoveBot()
}

func TestBotManager_RemoveBot_Empty(t *testing.T) {
	t.Parallel()
	bm := NewBotManager()

	removed := bm.RemoveBot()
	if removed != nil {
		t.Errorf("expected nil when no bots, got %v", removed)
	}
}

func TestBotManager_BotCount(t *testing.T) {
	t.Parallel()
	om := order.NewOrderManager()
	q := queue.NewOrderQueue()
	notifyNew := make(chan struct{}, 100)

	bm := NewBotManager()
	if bm.BotCount() != 0 {
		t.Errorf("expected 0 bots, got %d", bm.BotCount())
	}

	bm.AddBot(q, om, notifyNew)
	if bm.BotCount() != 1 {
		t.Errorf("expected 1 bot, got %d", bm.BotCount())
	}

	bm.AddBot(q, om, notifyNew)
	if bm.BotCount() != 2 {
		t.Errorf("expected 2 bots, got %d", bm.BotCount())
	}

	bm.RemoveBot()
	if bm.BotCount() != 1 {
		t.Errorf("expected 1 bot after removal, got %d", bm.BotCount())
	}

	bm.RemoveBot()
}

func TestBot_String(t *testing.T) {
	t.Parallel()
	om := order.NewOrderManager()
	q := queue.NewOrderQueue()
	notifyNew := make(chan struct{}, 100)

	bm := NewBotManager()
	b := bm.AddBot(q, om, notifyNew)

	time.Sleep(50 * time.Millisecond)

	s := b.String()
	if s != "Bot#1 (IDLE)" {
		t.Errorf("expected 'Bot#1 (IDLE)', got '%s'", s)
	}

	bm.RemoveBot()
}

func TestBot_StringProcessing(t *testing.T) {
	om := order.NewOrderManager()
	q := queue.NewOrderQueue()
	notifyNew := make(chan struct{}, 100)

	bm := NewBotManager()
	b := bm.AddBot(q, om, notifyNew)

	o := om.CreateOrder(order.Normal, "test")
	q.Enqueue(o)

	select {
	case notifyNew <- struct{}{}:
	default:
	}

	time.Sleep(100 * time.Millisecond)

	s := b.String()
	expected := "Bot#1 (PROCESSING: Order#1)"
	if s != expected {
		t.Errorf("expected '%s', got '%s'", expected, s)
	}

	bm.RemoveBot()
}

func TestBotManager_NotifyAll(t *testing.T) {
	om := order.NewOrderManager()
	q := queue.NewOrderQueue()
	notifyNew := make(chan struct{}, 100)

	bm := NewBotManager()
	bm.AddBot(q, om, notifyNew)
	bm.AddBot(q, om, notifyNew)

	o := om.CreateOrder(order.Normal, "test")
	q.Enqueue(o)

	bm.NotifyAll()

	time.Sleep(200 * time.Millisecond)

	if status, ok := om.GetOrderStatus(o.ID); !ok || status != order.Processing {
		t.Errorf("expected order to be processing after NotifyAll, got %d, ok=%v", status, ok)
	}

	bm.RemoveBot()
	bm.RemoveBot()
}

func TestBotManager_GetBots(t *testing.T) {
	t.Parallel()
	om := order.NewOrderManager()
	q := queue.NewOrderQueue()
	notifyNew := make(chan struct{}, 100)

	bm := NewBotManager()
	bm.AddBot(q, om, notifyNew)
	bm.AddBot(q, om, notifyNew)

	bots := bm.GetBots()
	if len(bots) != 2 {
		t.Errorf("expected 2 bots, got %d", len(bots))
	}

	bm.RemoveBot()
	bm.RemoveBot()
}
