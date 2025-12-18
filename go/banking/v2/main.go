package main

import (
	"fmt"
	"orderapi/transaction"
)

func main() {
	amounts := []int{500, 300, 200, 50, 80}
	types := []string{"item", "item", "item", "discount", "tax"}

	order := transaction.NewCalculateOrder()

	for i := 0; i < len(amounts); i++ {
		order.Cal(amounts[i], types[i])
	}

	result := order.Response()

	fmt.Printf("subtotal: %d,\ndiscount: %d,\ntax: %d,\ntotal: %d\n",
		result.Subtotal,
		result.Discount,
		result.Tax,
		result.Total)

	// Test with bonus features
	fmt.Println("\n--- With shipping and coupon ---")

	order2 := transaction.NewCalculateOrder()
	order2.Cal(1000, "item")
	order2.Cal(500, "item")
	order2.Cal(10, "coupon") // 10% off subtotal (150)
	order2.Cal(100, "tax")
	order2.Cal(50, "shipping")

	result2 := order2.Response()

	fmt.Printf("subtotal: %d,\ndiscount: %d,\ntax: %d,\nshipping: %d,\ntotal: %d\n",
		result2.Subtotal,
		result2.Discount,
		result2.Tax,
		result2.Shipping,
		result2.Total)
}
