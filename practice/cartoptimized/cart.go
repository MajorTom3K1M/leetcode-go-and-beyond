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
	Items        map[string]*Item
	Promotions   map[string]*Promotion
	BillDiscount int
	SubTotal     int
	Discount     int
}

func NewCart() *Cart {
	return &Cart{
		Items:        make(map[string]*Item),
		Promotions:   make(map[string]*Promotion),
		BillDiscount: 0,
		SubTotal:     0,
		Discount:     0,
	}
}

func (c *Cart) ApplyDiscount(itemName string, quantity int) {
	if promo, exists := c.Promotions[itemName]; exists {
		c.Discount += quantity * promo.Discount
	}
}

func (c *Cart) AddItem(item Item) {
	oldItem, ok := c.Items[item.Name]
	if ok {
		// remove old calculate sub total
		c.SubTotal -= oldItem.UnitPrice * oldItem.Quantity
		c.ApplyDiscount(item.Name, -oldItem.Quantity)

		oldItem.Quantity += item.Quantity
		oldItem.UnitPrice = item.UnitPrice

		c.SubTotal += item.UnitPrice * oldItem.Quantity
		c.ApplyDiscount(item.Name, oldItem.Quantity)

		return
	}

	c.SubTotal += item.UnitPrice * item.Quantity
	c.ApplyDiscount(item.Name, item.Quantity)

	c.Items[item.Name] = &item
}

func (c *Cart) AddCashPromotion(itemName string, discountPerUnit int) {
	item, hasItem := c.Items[itemName]

	if _, hasPromo := c.Promotions[itemName]; hasPromo && hasItem {
		c.ApplyDiscount(itemName, -item.Quantity)
	}

	c.Promotions[itemName] = &Promotion{
		Ref:      itemName,
		Discount: discountPerUnit,
	}

	if hasItem {
		c.ApplyDiscount(itemName, item.Quantity)
	}
}

func (c *Cart) AddSpecialPromotion(billDiscount int) {
	c.BillDiscount += billDiscount
}

func (c *Cart) TotalPrice() int {
	total := c.SubTotal - c.Discount - c.BillDiscount
	if total < 0 {
		return 0
	}
	return total
}
