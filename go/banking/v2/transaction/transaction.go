package transaction

import (
	"orderapi/calculate"
	"orderapi/service"
)

type CalculateOrder struct {
	subtotal int
	discount int
	tax      int
	shipping int
}

type Response struct {
	Subtotal int
	Discount int
	Tax      int
	Shipping int
	Total    int
}

func NewCalculateOrder() *CalculateOrder {
	return &CalculateOrder{
		subtotal: 0,
		discount: 0,
		tax:      0,
		shipping: 0,
	}
}

func (c *CalculateOrder) Cal(amount int, orderType string) bool {
	if !service.ValidateOrderType(orderType) {
		return false
	}

	if !service.ValidateAmount(amount) {
		return false
	}

	switch orderType {
	case "item":
		c.subtotal += amount
	case "discount":
		c.discount += amount
	case "tax":
		c.tax += amount
	case "shipping":
		c.shipping += amount
	case "coupon":
		percentDiscount := calculate.CalculatePercentage(c.subtotal, amount)
		c.discount += percentDiscount
	}

	return true
}

func (c *CalculateOrder) Response() *Response {
	total := calculate.CalculateTotal(c.subtotal, c.discount, c.tax, c.shipping)

	return &Response{
		Subtotal: c.subtotal,
		Discount: c.discount,
		Tax:      c.tax,
		Shipping: c.shipping,
		Total:    total,
	}
}
