package order

import (
	"fmt"
	"time"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusPreparing OrderStatus = "preparing"
	StatusReady     OrderStatus = "ready"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
	Name     string
	Price    int
	Quantity int
}

type Order struct {
	ID        string
	Items     []OrderItem
	Status    OrderStatus
	CreatedAt time.Time
}

type OrderManager struct {
	orders       map[string]*Order
	LastSequence int
}

var validTransitions = map[OrderStatus]map[OrderStatus]bool{
	StatusPending:   {StatusConfirmed: true, StatusCancelled: true},
	StatusConfirmed: {StatusPreparing: true, StatusCancelled: true},
	StatusPreparing: {StatusReady: true, StatusCancelled: true},
	StatusReady:     {StatusDelivered: true},
	StatusDelivered: {},
	StatusCancelled: {},
}

func NewOrderManager() *OrderManager {
	return &OrderManager{
		orders: map[string]*Order{},
	}
}

// CreateOrder creates a new order with pending status.
// Returns the order ID.
// Order ID format: "ORD-{sequential number starting from 1}"
// Example: "ORD-1", "ORD-2", etc.
func (om *OrderManager) CreateOrder(items []OrderItem) string {
	newSequence := om.LastSequence + 1
	om.LastSequence = newSequence
	newId := fmt.Sprintf("%s-%d", "ORD", newSequence)
	om.orders[newId] = &Order{
		ID:        newId,
		Items:     items,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
	return newId
}

// GetOrder returns an order by ID, or nil if not found.
func (om *OrderManager) GetOrder(id string) *Order {
	if order, ok := om.orders[id]; ok {
		return order
	}
	return nil
}

// UpdateStatus updates order status with validation.
// Valid transitions:
//
//	pending -> confirmed, cancelled
//	confirmed -> preparing, cancelled
//	preparing -> ready, cancelled
//	ready -> delivered
//	delivered -> (no transitions allowed)
//	cancelled -> (no transitions allowed)
//
// Returns true if transition was valid and applied.
func (om *OrderManager) UpdateStatus(id string, newStatus OrderStatus) bool {
	order, ok := om.orders[id]
	if !ok {
		return false
	}

	if validTransitions[order.Status][newStatus] {
		order.Status = newStatus
		return true
	}

	return false
}

// GetOrderTotal returns the total price for an order.
// Returns 0 if order not found.
func (om *OrderManager) GetOrderTotal(id string) int {
	if order, ok := om.orders[id]; ok {
		sum := 0
		for _, item := range order.Items {
			if item.Quantity < 0 {
				continue
			}
			sum += item.Price * item.Quantity
		}
		return sum
	}
	return 0
}

// GetOrdersByStatus returns all orders with the given status.
// Returns empty slice if none found.
func (om *OrderManager) GetOrdersByStatus(status OrderStatus) []*Order {
	result := []*Order{}
	for _, order := range om.orders {
		if order.Status == status {
			result = append(result, order)
		}
	}
	return result
}

// CancelOldPendingOrders cancels all pending orders older than the given duration.
// Returns the count of cancelled orders.
func (om *OrderManager) CancelOldPendingOrders(maxAge time.Duration) int {
	cancelledCount := 0
	for _, order := range om.orders {
		if order.Status == StatusPending && time.Since(order.CreatedAt) > maxAge {
			order.Status = StatusCancelled
			cancelledCount++
		}
	}
	return cancelledCount
}
