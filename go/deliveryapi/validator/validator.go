package validators

var ValidDeliveryType = map[string]bool{
	"standard": true,
	"express":  true,
}

var ValidPromoCode = map[string]bool{
	"DISCOUNT10":   true,
	"DRINK20":      true,
	"FREEDELIVERY": true,
	"FLAT50":       true,
}

var PromoType = map[string]string{
	"DISCOUNT10":   "percent_off",
	"DRINK20":      "percent_off",
	"FREEDELIVERY": "no_delivery_fee",
	"FLAT50":       "flat_off",
}

var PromoValue = map[string]int{
	"DISCOUNT10":   10,
	"DRINK20":      20,
	"FREEDELIVERY": 0,
	"FLAT50":       50,
}

var MinOrder = map[string]int{
	"FLAT50": 200,
}

func IsValidPromo(code string) bool {
	return ValidPromoCode[code]
}

func GetPromoType(code string) string {
	return PromoType[code]
}

func GetPromoValue(code string) int {
	return PromoValue[code]
}

func GetMinOrder(code string) int {
	return MinOrder[code]
}

func IsValidDeliveryType(priority string) bool {
	return ValidDeliveryType[priority]
}
