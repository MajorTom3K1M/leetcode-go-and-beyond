package cart

import (
	"cartapi/calculate"
	"cartapi/service"
)

type Cart struct {
	subtotal        int
	tax             int
	discount        int
	total           int
	itemCount       int
	taxableSubtotal int
}

type CheckoutResponse struct {
	ItemCount int
	Subtotal  int
	Discount  int
	Tax       int
	Total     int
}

func NewCart() *Cart {
	return &Cart{}
}

func (c *Cart) Add(price int, quantity int, itemType string) bool {
	if service.ValidateItemType(itemType) {
		amount := price * quantity
		c.subtotal += amount
		c.itemCount += quantity

		if service.IsTaxable(itemType) {
			c.taxableSubtotal += amount
		}

		return true
	} else if service.IsPromotion(itemType) {
		switch itemType {
		case "percent_off":
			discount := calculate.CalculatePercentOff(c.subtotal, price)
			c.discount += discount * quantity
		case "flat_off":
			c.discount += price * quantity
		}

		return true
	}

	return false
}

func (c *Cart) Checkout() *CheckoutResponse {
	taxableAfterDiscount := c.taxableSubtotal
	if c.subtotal > 0 && c.discount > 0 {
		taxableDiscount := (c.taxableSubtotal * c.discount) / c.subtotal
		taxableAfterDiscount = c.taxableSubtotal - taxableDiscount
	}

	c.tax = calculate.CalculateTax(taxableAfterDiscount)
	c.total = calculate.CalculateTotal(c.subtotal, c.discount, c.tax)

	return &CheckoutResponse{
		ItemCount: c.itemCount,
		Subtotal:  c.subtotal,
		Discount:  c.discount,
		Tax:       c.tax,
		Total:     c.total,
	}
}
