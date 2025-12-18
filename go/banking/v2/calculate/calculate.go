package calculate

func CalculateTotal(subtotal, discount, tax, shipping int) int {
	return subtotal - discount + tax + shipping
}

func CalculatePercentage(amount int, percent int) int {
	return (amount * percent) / 100
}
