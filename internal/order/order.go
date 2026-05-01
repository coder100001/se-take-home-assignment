package order

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type OrderType int

const (
	Normal OrderType = iota
	VIP
)

type OrderStatus int

const (
	Pending OrderStatus = iota
	Processing
	Complete
)

type Order struct {
	ID          int
	Type        OrderType
	Status      OrderStatus
	Description string
	CreatedAt   time.Time
}

var typeStrs = [...]string{"NORMAL", "VIP"}

var statusStrs = [...]string{"PENDING", "PROCESSING", "COMPLETE"}

func (o *Order) String() string {
	return fmt.Sprintf("Order#%d (%s, %s)", o.ID, typeStrs[o.Type], statusStrs[o.Status])
}

type OrderManager struct {
	nextID int64
	orders map[int]*Order
	mu     sync.RWMutex
}

func NewOrderManager() *OrderManager {
	return &OrderManager{
		orders: make(map[int]*Order),
	}
}

func (m *OrderManager) CreateOrder(orderType OrderType, description string) *Order {
	id := int(atomic.AddInt64(&m.nextID, 1))

	o := &Order{
		ID:          id,
		Type:        orderType,
		Status:      Pending,
		Description: description,
		CreatedAt:   time.Now(),
	}

	m.mu.Lock()
	m.orders[id] = o
	m.mu.Unlock()

	return o
}

func (m *OrderManager) GetOrder(id int) *Order {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.orders[id]
}

func (m *OrderManager) GetOrderStatus(id int) (OrderStatus, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if o, ok := m.orders[id]; ok {
		return o.Status, true
	}
	return Pending, false
}

func (m *OrderManager) UpdateStatus(id int, status OrderStatus) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if o, ok := m.orders[id]; ok {
		o.Status = status
	}
}

func (m *OrderManager) GetAllOrders() []Order {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]Order, 0, len(m.orders))
	for _, o := range m.orders {
		result = append(result, *o)
	}
	return result
}
