package main

import (
	"fmt"
	"sync"
	"time"
)

type Order struct {
	OrderID   string
	ProductID string
	Quantity  int
}

type OrderResult struct {
	OrderID string
	Success bool
	Message string
}

type Product struct {
	ProductID       string
	ItemName        string
	Quantity        int
	IsReserve       bool
	ReserveDuration time.Duration
	ReserveAt       time.Time
	ReserveAmount   int
}

type Inventory struct {
	orders   map[string]Order
	stock    map[string]Product
	orderSeq int
	mu       sync.RWMutex
}

func NewInventory() *Inventory {
	return &Inventory{
		orders:   make(map[string]Order),
		stock:    make(map[string]Product),
		orderSeq: 0,
	}
}

func (in *Inventory) AddStock(itemName string, qty int) {
	in.mu.Lock()
	defer in.mu.Unlock()

	if item, exists := in.stock[itemName]; exists {
		item.Quantity += qty
		in.stock[itemName] = item
		return
	}

	in.stock[itemName] = Product{
		ItemName:  itemName,
		ProductID: itemName,
		Quantity:  qty,
	}
}

func (in *Inventory) Purchase(itemName string, qty int) (bool, string) {
	in.mu.Lock()
	defer in.mu.Unlock()

	stock, exists := in.stock[itemName]
	if exists {
		stockQty := stock.Quantity
		if stock.IsReserve && time.Since(stock.ReserveAt) <= stock.ReserveDuration {
			stockQty -= stock.ReserveAmount
		} else if stock.IsReserve && time.Since(stock.ReserveAt) > stock.ReserveDuration {
			// reset
			stock.IsReserve = false
			stock.ReserveAmount = 0
			stock.ReserveAt = time.Time{}
			stock.ReserveDuration = 0
			in.stock[itemName] = stock
		}

		if qty > stockQty {
			return false, "insufficient stock"
		}

		stock.Quantity -= qty
		in.stock[itemName] = stock
	} else {
		return false, "product not found"
	}

	in.orderSeq++
	newID := fmt.Sprintf("%03d", in.orderSeq)
	newOrder := Order{
		OrderID:   newID,
		ProductID: itemName,
		Quantity:  qty,
	}
	in.orders[newID] = newOrder

	return true, fmt.Sprintf("stock now %d", stock.Quantity)
}

func (in *Inventory) GetStock(itemName string) int {
	in.mu.RLock()
	defer in.mu.RUnlock()

	return in.stock[itemName].Quantity
}

func (in *Inventory) ProcessOrders(orders []Order) []OrderResult {
	result := []OrderResult{}
	for _, order := range orders {
		success, message := in.Purchase(order.ProductID, order.Quantity)

		result = append(result, OrderResult{
			OrderID: order.OrderID,
			Success: success,
			Message: message,
		})
	}

	return result
}

func (in *Inventory) Reserve(itemName string, qty int, duration time.Duration) bool {
	in.mu.Lock()
	defer in.mu.Unlock()

	if item, exists := in.stock[itemName]; exists {
		if qty > item.Quantity {
			return false
		}

		item.IsReserve = true
		item.ReserveAmount = qty
		item.ReserveAt = time.Now()
		item.ReserveDuration = duration
		in.stock[itemName] = item

		return true
	}

	return false
}
