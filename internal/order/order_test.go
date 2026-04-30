package order

import (
	"testing"
)

func TestCreateOrder(t *testing.T) {
	t.Parallel()
	om := NewOrderManager()

	o := om.CreateOrder(Normal, "burger")
	if o.ID != 1 {
		t.Errorf("expected ID 1, got %d", o.ID)
	}
	if o.Type != Normal {
		t.Errorf("expected Normal type, got %d", o.Type)
	}
	if o.Status != Pending {
		t.Errorf("expected Pending status, got %d", o.Status)
	}
	if o.Description != "burger" {
		t.Errorf("expected description 'burger', got '%s'", o.Description)
	}
}

func TestCreateVIPOrder(t *testing.T) {
	t.Parallel()
	om := NewOrderManager()

	o := om.CreateOrder(VIP, "fries")
	if o.Type != VIP {
		t.Errorf("expected VIP type, got %d", o.Type)
	}
}

func TestOrderIDIncrement(t *testing.T) {
	t.Parallel()
	om := NewOrderManager()

	o1 := om.CreateOrder(Normal, "a")
	o2 := om.CreateOrder(VIP, "b")
	o3 := om.CreateOrder(Normal, "c")

	if o1.ID >= o2.ID {
		t.Errorf("expected o1.ID < o2.ID, got %d >= %d", o1.ID, o2.ID)
	}
	if o2.ID >= o3.ID {
		t.Errorf("expected o2.ID < o3.ID, got %d >= %d", o2.ID, o3.ID)
	}
}

func TestGetOrder(t *testing.T) {
	t.Parallel()
	om := NewOrderManager()

	o := om.CreateOrder(Normal, "test")
	got := om.GetOrder(o.ID)
	if got == nil {
		t.Fatal("expected order, got nil")
	}
	if got.ID != o.ID {
		t.Errorf("expected ID %d, got %d", o.ID, got.ID)
	}
}

func TestGetOrder_NotFound(t *testing.T) {
	t.Parallel()
	om := NewOrderManager()

	got := om.GetOrder(999)
	if got != nil {
		t.Errorf("expected nil for not found, got %v", got)
	}
}

func TestUpdateStatus(t *testing.T) {
	t.Parallel()
	om := NewOrderManager()

	o := om.CreateOrder(Normal, "test")
	om.UpdateStatus(o.ID, Processing)

	if status, ok := om.GetOrderStatus(o.ID); !ok || status != Processing {
		t.Errorf("expected Processing status, got %d, ok=%v", status, ok)
	}

	om.UpdateStatus(o.ID, Complete)
	if status, ok := om.GetOrderStatus(o.ID); !ok || status != Complete {
		t.Errorf("expected Complete status, got %d, ok=%v", status, ok)
	}
}

func TestGetOrderStatus(t *testing.T) {
	t.Parallel()
	om := NewOrderManager()

	status, ok := om.GetOrderStatus(999)
	if ok {
		t.Errorf("expected not found (ok=false), got ok=%v", ok)
	}
	if status != Pending {
		t.Errorf("expected Pending for not found, got %d", status)
	}
}

func TestGetAllOrders(t *testing.T) {
	t.Parallel()
	om := NewOrderManager()

	om.CreateOrder(Normal, "a")
	om.CreateOrder(VIP, "b")
	om.CreateOrder(Normal, "c")

	orders := om.GetAllOrders()
	if len(orders) != 3 {
		t.Errorf("expected 3 orders, got %d", len(orders))
	}
}

func TestOrderString(t *testing.T) {
	t.Parallel()

	o := &Order{ID: 1, Type: VIP, Status: Pending}
	expected := "Order#1 (VIP, PENDING)"
	if o.String() != expected {
		t.Errorf("expected '%s', got '%s'", expected, o.String())
	}

	o2 := &Order{ID: 2, Type: Normal, Status: Processing}
	expected2 := "Order#2 (NORMAL, PROCESSING)"
	if o2.String() != expected2 {
		t.Errorf("expected '%s', got '%s'", expected2, o2.String())
	}

	o3 := &Order{ID: 3, Type: Normal, Status: Complete}
	expected3 := "Order#3 (NORMAL, COMPLETE)"
	if o3.String() != expected3 {
		t.Errorf("expected '%s', got '%s'", expected3, o3.String())
	}
}
