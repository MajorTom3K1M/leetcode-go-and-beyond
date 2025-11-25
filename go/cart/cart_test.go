package cart

import "testing"

func TestEmptyCart(t *testing.T) {
	c := NewCart()
	if got := c.TotalPrice(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestSingleItem(t *testing.T) {
	c := NewCart()
	c.AddItem(Item{Name: "Bento", UnitPrice: 1000, Quantity: 2})
	if got := c.TotalPrice(); got != 2000 {
		t.Fatalf("expected 2000, got %d", got)
	}
}

func TestCashPromotion(t *testing.T) {
	c := NewCart()

	c.AddItem(Item{Name: "Chocolate", UnitPrice: 10, Quantity: 2})
	c.AddItem(Item{Name: "Pocky", UnitPrice: 20, Quantity: 2})
	c.AddItem(Item{Name: "Bento", UnitPrice: 15, Quantity: 2})

	c.AddCashPromotion("Bento", 5)

	if got := c.TotalPrice(); got != 80 {
		t.Fatalf("expected 60, got %d", got)
	}
}

func TestCashAndSpecialPromotion(t *testing.T) {
	c := NewCart()

	c.AddItem(Item{Name: "Chocolate", UnitPrice: 10, Quantity: 2})
	c.AddItem(Item{Name: "Pocky", UnitPrice: 20, Quantity: 2})
	c.AddItem(Item{Name: "Bento", UnitPrice: 15, Quantity: 2})

	c.AddCashPromotion("Bento", 5)
	c.AddSpecialPromotion(10)

	if got := c.TotalPrice(); got != 70 {
		t.Fatalf("expected 70, got %d", got)
	}
}

func TestTotalNeverNegative(t *testing.T) {
	c := NewCart()
	c.AddItem(Item{Name: "Candy", UnitPrice: 5, Quantity: 1})
	c.AddSpecialPromotion(1000)

	if got := c.TotalPrice(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}
