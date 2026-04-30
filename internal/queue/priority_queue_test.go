package queue

import (
	"sync"
	"testing"

	"github.com/feedme/order-controller/internal/order"
)

func TestEnqueueDequeue_Normal(t *testing.T) {
	t.Parallel()
	q := NewOrderQueue()
	om := order.NewOrderManager()

	o1 := om.CreateOrder(order.Normal, "test1")
	o2 := om.CreateOrder(order.Normal, "test2")
	q.Enqueue(o1)
	q.Enqueue(o2)

	got1 := q.Dequeue()
	got2 := q.Dequeue()

	if got1.ID != o1.ID {
		t.Errorf("expected first order ID %d, got %d", o1.ID, got1.ID)
	}
	if got2.ID != o2.ID {
		t.Errorf("expected second order ID %d, got %d", o2.ID, got2.ID)
	}
}

func TestEnqueueDequeue_VIP(t *testing.T) {
	t.Parallel()
	q := NewOrderQueue()
	om := order.NewOrderManager()

	o1 := om.CreateOrder(order.VIP, "vip1")
	o2 := om.CreateOrder(order.VIP, "vip2")
	q.Enqueue(o1)
	q.Enqueue(o2)

	got1 := q.Dequeue()
	got2 := q.Dequeue()

	if got1.ID != o1.ID {
		t.Errorf("expected first VIP order ID %d, got %d", o1.ID, got1.ID)
	}
	if got2.ID != o2.ID {
		t.Errorf("expected second VIP order ID %d, got %d", o2.ID, got2.ID)
	}
}

func TestDequeue_VIPPriority(t *testing.T) {
	t.Parallel()
	q := NewOrderQueue()
	om := order.NewOrderManager()

	normalOrder := om.CreateOrder(order.Normal, "normal1")
	vipOrder := om.CreateOrder(order.VIP, "vip1")
	q.Enqueue(normalOrder)
	q.Enqueue(vipOrder)

	got := q.Dequeue()
	if got.Type != order.VIP {
		t.Errorf("expected VIP order first, got type %d", got.Type)
	}
}

func TestDequeue_EmptyQueue(t *testing.T) {
	t.Parallel()
	q := NewOrderQueue()

	got := q.Dequeue()
	if got != nil {
		t.Errorf("expected nil from empty queue, got %v", got)
	}
}

func TestIsEmpty(t *testing.T) {
	t.Parallel()
	q := NewOrderQueue()
	om := order.NewOrderManager()

	if !q.IsEmpty() {
		t.Error("expected empty queue to be empty")
	}

	q.Enqueue(om.CreateOrder(order.Normal, "test"))
	if q.IsEmpty() {
		t.Error("expected non-empty queue to not be empty")
	}
}

func TestRequeueAt(t *testing.T) {
	t.Parallel()
	q := NewOrderQueue()
	om := order.NewOrderManager()

	o1 := om.CreateOrder(order.Normal, "n1")
	o2 := om.CreateOrder(order.Normal, "n2")
	o3 := om.CreateOrder(order.Normal, "n3")
	q.Enqueue(o1)
	q.Enqueue(o2)
	q.Enqueue(o3)

	q.Dequeue()
	q.RequeueAt(o1, 0)

	pending := q.PendingOrders()
	if len(pending) != 3 {
		t.Fatalf("expected 3 pending orders, got %d", len(pending))
	}
	if pending[0].ID != o1.ID {
		t.Errorf("expected requeued order at index 0, got ID %d", pending[0].ID)
	}
}

func TestRequeueAt_OutOfRange(t *testing.T) {
	t.Parallel()
	q := NewOrderQueue()
	om := order.NewOrderManager()

	o1 := om.CreateOrder(order.Normal, "n1")
	o2 := om.CreateOrder(order.Normal, "n2")
	q.Enqueue(o1)
	q.Enqueue(o2)

	q.Dequeue()
	q.RequeueAt(o1, 999)

	pending := q.PendingOrders()
	if len(pending) != 2 {
		t.Fatalf("expected 2 pending orders, got %d", len(pending))
	}
	if pending[len(pending)-1].ID != o1.ID {
		t.Errorf("expected requeued order at tail, got ID %d", pending[len(pending)-1].ID)
	}
}

func TestRemoveByID(t *testing.T) {
	t.Parallel()
	q := NewOrderQueue()
	om := order.NewOrderManager()

	o1 := om.CreateOrder(order.Normal, "n1")
	o2 := om.CreateOrder(order.Normal, "n2")
	q.Enqueue(o1)
	q.Enqueue(o2)

	removed, idx := q.RemoveByID(o1.ID)
	if removed == nil {
		t.Fatal("expected to find order, got nil")
	}
	if removed.ID != o1.ID {
		t.Errorf("expected removed order ID %d, got %d", o1.ID, removed.ID)
	}
	if idx != 0 {
		t.Errorf("expected index 0, got %d", idx)
	}

	pending := q.PendingOrders()
	if len(pending) != 1 {
		t.Fatalf("expected 1 pending order, got %d", len(pending))
	}
}

func TestRemoveByID_NotFound(t *testing.T) {
	t.Parallel()
	q := NewOrderQueue()

	removed, idx := q.RemoveByID(999)
	if removed != nil {
		t.Errorf("expected nil for not found, got %v", removed)
	}
	if idx != -1 {
		t.Errorf("expected index -1, got %d", idx)
	}
}

func TestPendingOrders(t *testing.T) {
	t.Parallel()
	q := NewOrderQueue()
	om := order.NewOrderManager()

	vip1 := om.CreateOrder(order.VIP, "v1")
	normal1 := om.CreateOrder(order.Normal, "n1")
	vip2 := om.CreateOrder(order.VIP, "v2")
	q.Enqueue(normal1)
	q.Enqueue(vip1)
	q.Enqueue(vip2)

	pending := q.PendingOrders()
	if len(pending) != 3 {
		t.Fatalf("expected 3 pending orders, got %d", len(pending))
	}
	if pending[0].Type != order.VIP || pending[1].Type != order.VIP {
		t.Error("expected VIP orders first in pending list")
	}
	if pending[2].Type != order.Normal {
		t.Error("expected Normal order last in pending list")
	}
}

func TestConcurrentAccess(t *testing.T) {
	t.Parallel()
	q := NewOrderQueue()
	om := order.NewOrderManager()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			o := om.CreateOrder(order.Normal, "test")
			q.Enqueue(o)
		}(i)
	}
	wg.Wait()

	if q.IsEmpty() {
		t.Error("expected queue to have orders after concurrent enqueues")
	}

	for i := 0; i < 100; i++ {
		o := q.Dequeue()
		if o == nil {
			t.Errorf("expected order %d, got nil", i)
		}
	}

	if !q.IsEmpty() {
		t.Error("expected queue to be empty after dequeuing all")
	}
}
