package controller

import (
	"testing"
	"time"

	"github.com/feedme/order-controller/internal/bot"
	"github.com/feedme/order-controller/internal/order"
)

func TestController_CreateNormalOrder(t *testing.T) {
	t.Parallel()
	c := NewController()

	o := c.CreateOrder(order.Normal, "burger")
	if o == nil {
		t.Fatal("expected order to be created")
	}
	if o.Type != order.Normal {
		t.Errorf("expected Normal order, got %d", o.Type)
	}

	pending, _, _, _ := c.GetStatus()
	found := false
	for _, p := range pending {
		if p.ID == o.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected order to be in pending list")
	}

	c.Shutdown()
}

func TestController_CreateVIPOrder(t *testing.T) {
	t.Parallel()
	c := NewController()

	o := c.CreateOrder(order.VIP, "fries")
	if o == nil {
		t.Fatal("expected order to be created")
	}
	if o.Type != order.VIP {
		t.Errorf("expected VIP order, got %d", o.Type)
	}

	c.Shutdown()
}

func TestController_OrderIDIncrement(t *testing.T) {
	t.Parallel()
	c := NewController()

	o1 := c.CreateOrder(order.Normal, "a")
	o2 := c.CreateOrder(order.VIP, "b")
	o3 := c.CreateOrder(order.Normal, "c")

	if o1.ID >= o2.ID {
		t.Errorf("expected o1.ID < o2.ID, got %d >= %d", o1.ID, o2.ID)
	}
	if o2.ID >= o3.ID {
		t.Errorf("expected o2.ID < o3.ID, got %d >= %d", o2.ID, o3.ID)
	}

	c.Shutdown()
}

func TestController_AddBot(t *testing.T) {
	t.Parallel()
	c := NewController()

	b := c.AddBot()
	if b == nil {
		t.Fatal("expected bot to be created")
	}
	if b.ID != 1 {
		t.Errorf("expected bot ID 1, got %d", b.ID)
	}

	c.Shutdown()
}

func TestController_RemoveBot(t *testing.T) {
	t.Parallel()
	c := NewController()

	c.AddBot()
	b2 := c.AddBot()

	removed := c.RemoveBot()
	if removed == nil {
		t.Fatal("expected bot to be removed")
	}
	if removed.ID != b2.ID {
		t.Errorf("expected newest bot (ID %d) to be removed, got ID %d", b2.ID, removed.ID)
	}

	c.Shutdown()
}

func TestController_RemoveBot_Empty(t *testing.T) {
	t.Parallel()
	c := NewController()

	removed := c.RemoveBot()
	if removed != nil {
		t.Errorf("expected nil when no bots, got %v", removed)
	}

	c.Shutdown()
}

func TestController_VIPPriority(t *testing.T) {
	c := NewController()

	normalOrder := c.CreateOrder(order.Normal, "normal1")
	vipOrder := c.CreateOrder(order.VIP, "vip1")

	c.AddBot()
	time.Sleep(200 * time.Millisecond)

	_, processing, _, bots := c.GetStatus()

	if len(bots) == 0 {
		t.Fatal("expected at least one bot")
	}

	if len(processing) > 0 {
		if processing[0].ID != vipOrder.ID {
			t.Errorf("expected VIP order (ID %d) to be processing first, got ID %d", vipOrder.ID, processing[0].ID)
		}
	} else {
		for _, b := range bots {
			if b.Status == bot.BotProcessing && b.Order != nil {
				if b.Order.ID != vipOrder.ID {
					t.Errorf("expected VIP order (ID %d) to be processing first, got ID %d", vipOrder.ID, b.Order.ID)
				}
			}
		}
	}

	_ = normalOrder
	c.Shutdown()
}

func TestController_BotProcessesOrder(t *testing.T) {
	c := NewController()

	o := c.CreateOrder(order.Normal, "test")
	c.AddBot()

	time.Sleep(200 * time.Millisecond)

	_, processing, _, _ := c.GetStatus()
	if len(processing) == 0 {
		_, _, complete, _ := c.GetStatus()
		if len(complete) == 0 {
			t.Error("expected order to be processing or complete")
		}
	}

	_ = o
	c.Shutdown()
}

func TestController_RemoveBotWhileProcessing(t *testing.T) {
	c := NewController()

	o := c.CreateOrder(order.Normal, "test")
	c.AddBot()

	time.Sleep(200 * time.Millisecond)

	c.RemoveBot()
	time.Sleep(100 * time.Millisecond)

	pending, _, _, _ := c.GetStatus()
	found := false
	for _, p := range pending {
		if p.ID == o.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected order to be back in pending after bot removal")
	}

	c.Shutdown()
}

func TestController_Shutdown(t *testing.T) {
	c := NewController()

	c.AddBot()
	c.AddBot()
	c.CreateOrder(order.Normal, "test")

	c.Shutdown()

	_, _, _, bots := c.GetStatus()
	if len(bots) != 0 {
		t.Errorf("expected 0 bots after shutdown, got %d", len(bots))
	}
}
