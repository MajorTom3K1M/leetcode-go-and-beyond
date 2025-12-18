package service

var ItemType = map[string]bool{
	"product": true,
	"gift":    true,
}

var TaxableItemType = map[string]bool{
	"product": true,
}

var PromotionType = map[string]bool{
	"percent_off": true,
	"flat_off":    true,
}

func ValidateItemType(itemType string) bool {
	return ItemType[itemType]
}

func IsTaxable(itemType string) bool {
	return TaxableItemType[itemType]
}

func IsPromotion(itemType string) bool {
	return PromotionType[itemType]
}
