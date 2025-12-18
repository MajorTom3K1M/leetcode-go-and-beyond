package main

import (
	"cartapi/cart"
	"fmt"
)

func main() {
	testcases := []struct {
		price    int
		quantity int
		itemType string
	}{
		{150, 2, "product"},
		{250, 1, "product"},
		{100, 3, "product"},
		{300, 1, "product"},
	}

	shoppingCart := cart.NewCart()

	for i := 0; i < len(testcases); i++ {
		shoppingCart.Add(testcases[i].price, testcases[i].quantity, testcases[i].itemType)
	}

	result := shoppingCart.Checkout()

	fmt.Printf("items: %d,\nsubtotal: %d,\ndiscount: %d,\ntax: %d,\ntotal: %d\n",
		result.ItemCount,
		result.Subtotal,
		result.Discount,
		result.Tax,
		result.Total)

	// Scenario 2: With promotions
	fmt.Println("\n--- With Promotions ---")

	cart2 := cart.NewCart()
	cart2.Add(500, 2, "product")    // 1000
	cart2.Add(200, 3, "product")    // 600
	cart2.Add(100, 1, "gift")       // 100 (no tax on gifts)
	cart2.Add(20, 1, "percent_off") // 20% discount
	cart2.Add(50, 1, "flat_off")    // 50 baht off

	result2 := cart2.Checkout()

	fmt.Printf("items: %d,\nsubtotal: %d,\ndiscount: %d,\ntax: %d,\ntotal: %d\n",
		result2.ItemCount,
		result2.Subtotal,
		result2.Discount,
		result2.Tax,
		result2.Total)
}
