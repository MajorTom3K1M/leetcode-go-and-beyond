package order

import (
	"deliveryapi/models"
	"deliveryapi/pricing"
	validators "deliveryapi/validator"
)

type Order struct {
	items        map[string]*models.OrderItem
	delivery     models.DeliveryInfo
	promoCode    string
	foodTotal    int
	drinkTotal   int
	dessertTotal int
	subTotal     int
}

func NewOrder() *Order {
	return &Order{
		items: make(map[string]*models.OrderItem),
	}
}

func (o *Order) CalculateByCategory(price int, quantity int, category string) {
	switch category {
	case "food":
		o.foodTotal += price * quantity
	case "drink":
		o.drinkTotal += price * quantity
	case "dessert":
		o.dessertTotal += price * quantity
	}
	o.subTotal += price * quantity
}

func (o *Order) AddItem(name string, price int, quantity int, category string) bool {
	if _, exists := o.items[name]; exists {
		return false // already exists
	}

	if price <= 0 || quantity <= 0 {
		return false // insufficient price and quantity
	}

	o.CalculateByCategory(price, quantity, category)
	newItem := &models.OrderItem{
		Name:     name,
		Price:    price,
		Quantity: quantity,
		Category: category,
	}
	o.items[name] = newItem

	return true
}

func (o *Order) SetDelivery(distance int, priority string) bool {
	if !validators.IsValidDeliveryType(priority) || distance <= 0 {
		return false
	}

	o.delivery = models.DeliveryInfo{
		Distance: distance,
		Priority: priority,
	}
	return true
}

func (o *Order) ApplyPromo(code string) bool {
	if !validators.IsValidPromo(code) {
		return false
	}
	o.promoCode = code
	return true
}

func (o *Order) CalculateDiscount() (int, string) {
	if !validators.IsValidPromo(o.promoCode) {
		return 0, ""
	}
	promoType := validators.GetPromoType(o.promoCode)
	promoValue := validators.GetPromoValue(o.promoCode)
	minOrder := validators.GetMinOrder(o.promoCode)
	if minOrder > 0 && o.subTotal < minOrder {
		return 0, ""
	}
	discount := 0
	switch promoType {
	case "percent_off":
		if o.promoCode == "DISCOUNT10" {
			discount = (o.foodTotal * promoValue) / 100
		} else if o.promoCode == "DRINK20" {
			discount = (o.drinkTotal * promoValue) / 100
		}
	case "flat_off":
		discount = promoValue
	case "no_delivery_fee":
		discount = 0
	}

	return discount, promoType
}

func (o *Order) Checkout() *models.OrderSummary {
	items := make([]models.OrderItem, 0, len(o.items))

	for _, item := range o.items {
		items = append(items, *item)
	}

	subtotal := o.subTotal
	deliveryInfo := o.delivery
	deliveryFee := pricing.CalculateDeliveryFee(deliveryInfo.Distance, deliveryInfo.Priority)
	smallOrderFee := pricing.CalculateSmallOrderFee(subtotal)

	discount, promoType := o.CalculateDiscount()

	if promoType == "no_delivery_fee" {
		deliveryFee = 0
	}

	serviceFee := pricing.CalculateServiceFee(subtotal - discount + deliveryFee + smallOrderFee)
	total := subtotal - discount + deliveryFee + smallOrderFee + serviceFee

	estimatedTime := pricing.CalculateEstimatedTime(deliveryInfo.Distance, deliveryInfo.Priority)

	summary := &models.OrderSummary{
		Items:         items,
		FoodTotal:     o.foodTotal,
		DrinkTotal:    o.drinkTotal,
		DessertTotal:  o.dessertTotal,
		Subtotal:      o.subTotal,
		DeliveryFee:   deliveryFee,
		SmallOrderFee: smallOrderFee,
		Discount:      discount,
		ServiceFee:    serviceFee,
		Total:         total,
		EstimatedTime: estimatedTime,
	}
	return summary
}
