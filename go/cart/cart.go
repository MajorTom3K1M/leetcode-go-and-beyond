package cart

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
	items             []Item
	promotions        map[string]Promotion
	specialPromotions []SpecialPromotion
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
	c.items = append(c.items, item)
}

// Cash promotion: discount per unit for a given item name.
// Example: AddCashPromotion("Bento", 5) = 5 off per Bento.
func (c *Cart) AddCashPromotion(itemName string, discountPerUnit int) {
	if discountPerUnit < 0 {
		discountPerUnit = 0
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
	c.specialPromotions = append(c.specialPromotions, SpecialPromotion{
		Discount: discount,
	})
}

// TotalPrice calculates the total after all promotions.
func (c *Cart) TotalPrice() int {
	totalPrice := 0
	for _, item := range c.items {
		discount, ok := c.promotions[item.Name]
		totalPrice += item.UnitPrice * item.Quantity
		if ok {
			totalPrice -= discount.Discount * item.Quantity
		}
	}

	for _, d := range c.specialPromotions {
		totalPrice -= d.Discount
	}

	if totalPrice < 0 {
		return 0
	}

	return totalPrice
}
