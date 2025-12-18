package calculate

var TaxValue int = 7

func CalculateTax(amount int) int {
	return amount * TaxValue / 100
}

func CalculatePercentOff(amount int, percent int) int {
	return amount * percent / 100
}

func CalculateTotal(subtotal, discount, tax int) int {
	return subtotal + tax - discount
}
