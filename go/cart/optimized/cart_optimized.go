package cart_optimized

type Item struct {
	Name      string
	UnitPrice int
	Quantity  int
}

type Promotion struct {
	Ref      string
	Discount int
}

type SpecialPromotion struct {
	Discount int
}

type Cart struct {
	items              []Item
	promotions         map[string]Promotion
	specialPromotions  []SpecialPromotion
	totalPriceCache    int
	specialDiscountSum int
}

func NewCart() *Cart {
	return &Cart{
		items:      []Item{},
		promotions: map[string]Promotion{},
	}
}

func (c *Cart) AddItem(item Item) {
	if item.Quantity <= 0 {
		item.Quantity = 1
	}

	promotion, ok := c.promotions[item.Name]
	if ok {
		discountPrice := max(item.UnitPrice-promotion.Discount, 0)
		c.totalPriceCache += discountPrice * item.Quantity
	} else {
		c.totalPriceCache += item.UnitPrice * item.Quantity
	}

	c.items = append(c.items, item)
}

// Cash promotion: discount per unit for a given item name.
// Example: AddCashPromotion("Bento", 5) = 5 off per Bento.
func (c *Cart) AddCashPromotion(itemName string, discountPerUnit int) {
	if discountPerUnit < 0 {
		discountPerUnit = 0
	}

	// if there old promotion remove it effect first
	if oldPromotion, ok := c.promotions[itemName]; ok {
		for _, item := range c.items {
			if item.Name == itemName {
				c.totalPriceCache += oldPromotion.Discount * item.Quantity
			}
		}
	}

	// apply new promotion effect
	for _, item := range c.items {
		if item.Name == itemName {
			c.totalPriceCache -= discountPerUnit * item.Quantity
		}
	}

	c.promotions[itemName] = Promotion{
		Ref:      itemName,
		Discount: discountPerUnit,
	}
}

// Special promotion: fixed discount on whole order.
// Example: AddSpecialPromotion(10) = order - 10.
func (c *Cart) AddSpecialPromotion(discount int) {
	if discount <= 0 {
		return
	}
	c.specialDiscountSum += discount
	c.specialPromotions = append(c.specialPromotions, SpecialPromotion{
		Discount: discount,
	})
}

// TotalPrice calculates the total after all promotions.
func (c *Cart) TotalPrice() int {
	totalPrice := c.totalPriceCache - c.specialDiscountSum

	if totalPrice < 0 {
		return 0
	}

	return totalPrice
}
