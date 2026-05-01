package queue

import (
	"github.com/feedme/order-controller/internal/order"
	"sync"
)

type OrderQueue struct {
	vipOrders    []*order.Order
	normalOrders []*order.Order
	mu           sync.Mutex
}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{
		vipOrders:    make([]*order.Order, 0),
		normalOrders: make([]*order.Order, 0),
	}
}

func (q *OrderQueue) Enqueue(o *order.Order) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if o.Type == order.VIP {
		q.vipOrders = append(q.vipOrders, o)
	} else {
		q.normalOrders = append(q.normalOrders, o)
	}
}

func (q *OrderQueue) Dequeue() *order.Order {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.vipOrders) > 0 {
		o := q.vipOrders[0]
		q.vipOrders = q.vipOrders[1:]
		return o
	}

	if len(q.normalOrders) > 0 {
		o := q.normalOrders[0]
		q.normalOrders = q.normalOrders[1:]
		return o
	}

	return nil
}

func (q *OrderQueue) RequeueAt(o *order.Order, index int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if o.Type == order.VIP {
		if index < 0 {
			index = 0
		}
		if index > len(q.vipOrders) {
			index = len(q.vipOrders)
		}
		q.vipOrders = append(q.vipOrders, nil)
		copy(q.vipOrders[index+1:], q.vipOrders[index:])
		q.vipOrders[index] = o
	} else {
		if index < 0 {
			index = 0
		}
		if index > len(q.normalOrders) {
			index = len(q.normalOrders)
		}
		q.normalOrders = append(q.normalOrders, nil)
		copy(q.normalOrders[index+1:], q.normalOrders[index:])
		q.normalOrders[index] = o
	}
}

func (q *OrderQueue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.vipOrders) == 0 && len(q.normalOrders) == 0
}

func (q *OrderQueue) PendingOrders() []*order.Order {
	q.mu.Lock()
	defer q.mu.Unlock()

	result := make([]*order.Order, 0, len(q.vipOrders)+len(q.normalOrders))
	result = append(result, q.vipOrders...)
	result = append(result, q.normalOrders...)
	return result
}

func (q *OrderQueue) RemoveByID(id int) (*order.Order, int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for i, o := range q.vipOrders {
		if o.ID == id {
			q.vipOrders = append(q.vipOrders[:i], q.vipOrders[i+1:]...)
			return o, i
		}
	}

	for i, o := range q.normalOrders {
		if o.ID == id {
			q.normalOrders = append(q.normalOrders[:i], q.normalOrders[i+1:]...)
			return o, i
		}
	}

	return nil, -1
}
