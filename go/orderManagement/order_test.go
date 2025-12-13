package order

import (
	"testing"
	"time"
)

func TestCreateAndGetOrder(t *testing.T) {
	om := NewOrderManager()

	items := []OrderItem{
		{Name: "Burger", Price: 100, Quantity: 2},
		{Name: "Fries", Price: 50, Quantity: 1},
	}

	id := om.CreateOrder(items)
	if id != "ORD-1" {
		t.Errorf("expected ORD-1, got %s", id)
	}

	order := om.GetOrder(id)
	if order == nil {
		t.Fatal("order should not be nil")
	}
	if order.Status != StatusPending {
		t.Errorf("expected pending status, got %s", order.Status)
	}
	if len(order.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(order.Items))
	}

	// Second order
	id2 := om.CreateOrder([]OrderItem{{Name: "Pizza", Price: 200, Quantity: 1}})
	if id2 != "ORD-2" {
		t.Errorf("expected ORD-2, got %s", id2)
	}
}

func TestGetOrderNotFound(t *testing.T) {
	om := NewOrderManager()
	order := om.GetOrder("ORD-999")
	if order != nil {
		t.Error("expected nil for non-existent order")
	}
}

func TestGetOrderTotal(t *testing.T) {
	om := NewOrderManager()

	items := []OrderItem{
		{Name: "Burger", Price: 100, Quantity: 2}, // 200
		{Name: "Fries", Price: 50, Quantity: 3},   // 150
	}
	id := om.CreateOrder(items)

	total := om.GetOrderTotal(id)
	if total != 350 {
		t.Errorf("expected 350, got %d", total)
	}

	// Non-existent order
	if om.GetOrderTotal("ORD-999") != 0 {
		t.Error("expected 0 for non-existent order")
	}
}

func TestValidStatusTransitions(t *testing.T) {
	om := NewOrderManager()
	id := om.CreateOrder([]OrderItem{{Name: "Soup", Price: 80, Quantity: 1}})

	// pending -> confirmed
	if !om.UpdateStatus(id, StatusConfirmed) {
		t.Error("pending -> confirmed should be valid")
	}

	// confirmed -> preparing
	if !om.UpdateStatus(id, StatusPreparing) {
		t.Error("confirmed -> preparing should be valid")
	}

	// preparing -> ready
	if !om.UpdateStatus(id, StatusReady) {
		t.Error("preparing -> ready should be valid")
	}

	// ready -> delivered
	if !om.UpdateStatus(id, StatusDelivered) {
		t.Error("ready -> delivered should be valid")
	}

	// Verify final status
	order := om.GetOrder(id)
	if order.Status != StatusDelivered {
		t.Errorf("expected delivered, got %s", order.Status)
	}
}

func TestInvalidStatusTransitions(t *testing.T) {
	om := NewOrderManager()

	// Test: pending -> ready (invalid, should skip confirmed/preparing)
	id1 := om.CreateOrder([]OrderItem{{Name: "Tea", Price: 30, Quantity: 1}})
	if om.UpdateStatus(id1, StatusReady) {
		t.Error("pending -> ready should be invalid")
	}

	// Test: delivered -> cancelled (invalid)
	id2 := om.CreateOrder([]OrderItem{{Name: "Coffee", Price: 40, Quantity: 1}})
	om.UpdateStatus(id2, StatusConfirmed)
	om.UpdateStatus(id2, StatusPreparing)
	om.UpdateStatus(id2, StatusReady)
	om.UpdateStatus(id2, StatusDelivered)
	if om.UpdateStatus(id2, StatusCancelled) {
		t.Error("delivered -> cancelled should be invalid")
	}

	// Test: cancelled -> confirmed (invalid)
	id3 := om.CreateOrder([]OrderItem{{Name: "Juice", Price: 50, Quantity: 1}})
	om.UpdateStatus(id3, StatusCancelled)
	if om.UpdateStatus(id3, StatusConfirmed) {
		t.Error("cancelled -> confirmed should be invalid")
	}

	// Test: non-existent order
	if om.UpdateStatus("ORD-999", StatusConfirmed) {
		t.Error("updating non-existent order should return false")
	}
}

func TestGetOrdersByStatus(t *testing.T) {
	om := NewOrderManager()

	om.CreateOrder([]OrderItem{{Name: "A", Price: 10, Quantity: 1}}) // pending
	om.CreateOrder([]OrderItem{{Name: "B", Price: 20, Quantity: 1}}) // pending
	id3 := om.CreateOrder([]OrderItem{{Name: "C", Price: 30, Quantity: 1}})
	om.UpdateStatus(id3, StatusConfirmed) // confirmed

	pending := om.GetOrdersByStatus(StatusPending)
	if len(pending) != 2 {
		t.Errorf("expected 2 pending orders, got %d", len(pending))
	}

	confirmed := om.GetOrdersByStatus(StatusConfirmed)
	if len(confirmed) != 1 {
		t.Errorf("expected 1 confirmed order, got %d", len(confirmed))
	}

	ready := om.GetOrdersByStatus(StatusReady)
	if len(ready) != 0 {
		t.Errorf("expected 0 ready orders, got %d", len(ready))
	}
}

func TestCancelOldPendingOrders(t *testing.T) {
	om := NewOrderManager()

	// Create orders
	id1 := om.CreateOrder([]OrderItem{{Name: "Old1", Price: 10, Quantity: 1}})
	id2 := om.CreateOrder([]OrderItem{{Name: "Old2", Price: 20, Quantity: 1}})
	id3 := om.CreateOrder([]OrderItem{{Name: "New", Price: 30, Quantity: 1}})

	// Manually set old timestamps for testing
	om.GetOrder(id1).CreatedAt = time.Now().Add(-2 * time.Hour)
	om.GetOrder(id2).CreatedAt = time.Now().Add(-90 * time.Minute)
	// id3 stays recent

	// Also make one confirmed (should not be cancelled even if old)
	id4 := om.CreateOrder([]OrderItem{{Name: "OldConfirmed", Price: 40, Quantity: 1}})
	om.GetOrder(id4).CreatedAt = time.Now().Add(-3 * time.Hour)
	om.UpdateStatus(id4, StatusConfirmed)

	// Cancel orders older than 1 hour
	count := om.CancelOldPendingOrders(1 * time.Hour)

	if count != 2 {
		t.Errorf("expected 2 cancelled, got %d", count)
	}

	if om.GetOrder(id1).Status != StatusCancelled {
		t.Error("id1 should be cancelled")
	}
	if om.GetOrder(id2).Status != StatusCancelled {
		t.Error("id2 should be cancelled")
	}
	if om.GetOrder(id3).Status != StatusPending {
		t.Error("id3 should still be pending")
	}
	if om.GetOrder(id4).Status != StatusConfirmed {
		t.Error("id4 should still be confirmed")
	}
}
