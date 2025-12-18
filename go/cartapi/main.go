package main

import (
	"cartapi/cart"
	"fmt"
)

func main() {
	items := []int{150, 250, 100, 300} // prices
	quantities := []int{2, 1, 3, 1}    // quantities
	types := []string{"product", "product", "product", "product"}

	shoppingCart := cart.NewCart()

	for i := 0; i < len(items); i++ {
		shoppingCart.Add(items[i], quantities[i], types[i])
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
