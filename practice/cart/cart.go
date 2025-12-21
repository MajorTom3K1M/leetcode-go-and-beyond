package main

type Item struct {
	Name      string
	UnitPrice int
	Quantity  int
}

type Promotion struct {
	Ref      string
	Discount int
}

type Cart struct {
	Items      map[string]*Item
	Promotions map[string]*Promotion
}

func NewCart() *Cart {
	return &Cart{
		Items:      make(map[string]*Item),
		Promotions: make(map[string]*Promotion),
	}
}

func (c *Cart) AddItem(item Item) {
	if item.Quantity <= 0 {
		item.Quantity = 1
	}

	if item.UnitPrice <= 0 {
		return
	}

	if oldItem, ok := c.Items[item.Name]; ok {
		oldItem.Quantity += item.Quantity
		oldItem.UnitPrice = item.UnitPrice
		return
	}

	c.Items[item.Name] = &item
}

func (c *Cart) AddPromotion(itemName string, dicountPerUnit int) {
	promotion, ok := c.Promotions[itemName]
	if ok {
		promotion.Discount = dicountPerUnit
		return
	}

	c.Promotions[itemName] = &Promotion{
		Ref:      itemName,
		Discount: dicountPerUnit,
	}
}

func (c *Cart) TotalPrice() int {
	total := 0
	for _, item := range c.Items {
		discount := 0
		if promotion, ok := c.Promotions[item.Name]; ok {
			discount += promotion.Discount * item.Quantity
		}
		includeDisount := (item.Quantity * item.UnitPrice) - discount
		if includeDisount > 0 {
			total += includeDisount
		}
	}

	return total
}
