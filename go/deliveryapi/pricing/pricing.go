package pricing

var DeliveryFeeType = map[string]int{
	"standard": 10,
	"express":  15,
}

var DeliveryEstimatedTime = map[string]int{
	"standard": 15,
	"express":  10,
}

func CalculateDeliveryFee(distance int, priority string) int {
	return distance * DeliveryFeeType[priority]
}

func CalculateSmallOrderFee(subtotal int) int {
	if subtotal < 100 {
		return 20
	}
	return 0
}

func CalculateServiceFee(amount int) int {
	return amount * 5 / 100
}

func CalculateEstimatedTime(distance int, priority string) int {
	val, exists := DeliveryEstimatedTime[priority]
	if !exists {
		return 0
	}

	if priority == "standard" {
		return val + distance*3
	} else if priority == "express" {
		return val + distance*1
	}

	return 0
}
