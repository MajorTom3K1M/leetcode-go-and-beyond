package models

type OrderItem struct {
	Name     string
	Price    int
	Quantity int
	Category string // "food", "drink", "dessert"
}

type DeliveryInfo struct {
	Distance int    // in kilometers
	Priority string // "standard", "express"
}

type OrderSummary struct {
	Items         []OrderItem
	FoodTotal     int
	DrinkTotal    int
	DessertTotal  int
	Subtotal      int
	DeliveryFee   int
	SmallOrderFee int // Extra fee if subtotal < 100
	Discount      int
	ServiceFee    int // 5% platform fee
	Total         int
	EstimatedTime int // minutes
}
