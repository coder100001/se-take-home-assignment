package controller

import (
	"github.com/feedme/order-controller/internal/bot"
	"github.com/feedme/order-controller/internal/order"
	"github.com/feedme/order-controller/internal/queue"
	"sort"
)

type botInfo struct {
	ID     int
	Status bot.BotStatus
	Order  *order.Order
}

type Controller struct {
	orderMgr  *order.OrderManager
	queue     *queue.OrderQueue
	botMgr    *bot.BotManager
	notifyNew chan struct{}
}

func NewController() *Controller {
	return &Controller{
		orderMgr:  order.NewOrderManager(),
		queue:     queue.NewOrderQueue(),
		botMgr:    bot.NewBotManager(),
		notifyNew: make(chan struct{}, 100),
	}
}

func (c *Controller) CreateOrder(orderType order.OrderType, description string) *order.Order {
	o := c.orderMgr.CreateOrder(orderType, description)
	c.queue.Enqueue(o)
	c.botMgr.NotifyAll()
	return o
}

func (c *Controller) AddBot() *bot.Bot {
	b := c.botMgr.AddBot(c.queue, c.orderMgr, c.notifyNew)
	select {
	case c.notifyNew <- struct{}{}:
	default:
	}
	return b
}

func (c *Controller) RemoveBot() *bot.Bot {
	return c.botMgr.RemoveBot()
}

func (c *Controller) GetStatus() (pending []order.Order, processing []order.Order, complete []order.Order, bots []botInfo) {
	orders := c.orderMgr.GetAllOrders()
	for _, o := range orders {
		switch o.Status {
		case order.Pending:
			pending = append(pending, o)
		case order.Processing:
			processing = append(processing, o)
		case order.Complete:
			complete = append(complete, o)
		}
	}

	sort.Slice(pending, func(i, j int) bool {
		return pending[i].Type > pending[j].Type
	})

	botList := c.botMgr.GetBots()
	bots = make([]botInfo, 0, len(botList))
	for _, b := range botList {
		status, curOrder := b.GetStatus()
		bots = append(bots, botInfo{
			ID:     b.ID,
			Status: status,
			Order:  curOrder,
		})
	}

	return pending, processing, complete, bots
}

func (c *Controller) Shutdown() {
	for c.botMgr.BotCount() > 0 {
		c.botMgr.RemoveBot()
	}
}
