package service

var validOrderTypes = map[string]bool{
	"item":     true,
	"discount": true,
	"tax":      true,
	"shipping": true,
	"coupon":   true,
}

func ValidateOrderType(orderType string) bool {
	return validOrderTypes[orderType]
}

func ValidateAmount(amount int) bool {
	return amount >= 0
}
