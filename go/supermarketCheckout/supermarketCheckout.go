package main

type Item struct {
	Name      string
	UnitPrice int
	Quantity  int
}

type BuyXGetYFree struct {
	Name string
	X    int
	Y    int
}

type Checkout struct {
	items       []Item
	bxgyPromos  map[string]BuyXGetYFree
	couponTotal int
}

func NewCheckout() *Checkout {
	return &Checkout{
		items:      make([]Item, 0),
		bxgyPromos: make(map[string]BuyXGetYFree),
	}
}

func (c *Checkout) AddItem(item Item) {
	if item.Quantity <= 0 {
		item.Quantity = 1
	}

	if item.UnitPrice < 0 {
		item.UnitPrice = 0
	}

	c.items = append(c.items, item)
}

func (c *Checkout) AddBuyXGetYFree(promo BuyXGetYFree) {
	if promo.X <= 0 || promo.Y <= 0 || promo.Name == "" {
		return
	}
	c.bxgyPromos[promo.Name] = promo
}

func (c *Checkout) AddCoupon(discount int) {
	if discount <= 0 {
		return
	}
	c.couponTotal += discount
}

func (c *Checkout) Total() int {
	subtotal := 0

	for _, item := range c.items {
		qty := item.Quantity
		unitPrice := item.UnitPrice

		if promo, ok := c.bxgyPromos[item.Name]; ok {
			groupSize := promo.X + promo.Y
			if groupSize > 0 {
				fullGroups := qty / groupSize
				freeUnits := fullGroups * promo.Y
				payableUnits := qty - freeUnits
				if payableUnits < 0 {
					payableUnits = 0
				}
				subtotal += payableUnits * unitPrice
			} else {
				subtotal += qty * unitPrice
			}
		} else {
			subtotal += qty * unitPrice
		}
	}

	total := subtotal - c.couponTotal
	if total < 0 {
		return 0
	}
	return total
}
