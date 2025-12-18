package main

import (
	"deliveryapi/models"
	"deliveryapi/order"
	"fmt"
)

func main() {
	// Scenario 1: Normal order
	fmt.Println("=== Scenario 1: Normal Order ===")

	order1 := order.NewOrder()

	order1.AddItem("Pad Thai", 89, 2, "food")
	order1.AddItem("Green Curry", 120, 1, "food")
	order1.AddItem("Thai Tea", 45, 2, "drink")
	order1.AddItem("Mango Sticky Rice", 79, 1, "dessert")

	order1.SetDelivery(3, "standard") // 3km, standard delivery
	order1.ApplyPromo("DISCOUNT10")   // 10% off food items

	summary1 := order1.Checkout()
	printSummary(summary1)

	// Scenario 2: Small order with express delivery
	fmt.Println("\n=== Scenario 2: Small Order + Express ===")

	order2 := order.NewOrder()

	order2.AddItem("Spring Rolls", 59, 1, "food")
	order2.AddItem("Coke", 25, 1, "drink")

	order2.SetDelivery(5, "express")  // 5km, express delivery
	order2.ApplyPromo("FREEDELIVERY") // Free delivery

	summary2 := order2.Checkout()
	printSummary(summary2)

	// Scenario 3: Drinks only (special case)
	fmt.Println("\n=== Scenario 3: Drinks Only ===")

	order3 := order.NewOrder()

	order3.AddItem("Thai Tea", 45, 3, "drink")
	order3.AddItem("Coffee", 60, 2, "drink")

	order3.SetDelivery(2, "standard")
	order3.ApplyPromo("DRINK20") // 20% off drinks

	summary3 := order3.Checkout()
	printSummary(summary3)
}

func printSummary(s *models.OrderSummary) {
	fmt.Printf("Food Total:      %d\n", s.FoodTotal)
	fmt.Printf("Drink Total:     %d\n", s.DrinkTotal)
	fmt.Printf("Dessert Total:   %d\n", s.DessertTotal)
	fmt.Printf("Subtotal:        %d\n", s.Subtotal)
	fmt.Printf("Delivery Fee:    %d\n", s.DeliveryFee)
	fmt.Printf("Small Order Fee: %d\n", s.SmallOrderFee)
	fmt.Printf("Discount:        -%d\n", s.Discount)
	fmt.Printf("Service Fee:     %d\n", s.ServiceFee)
	fmt.Printf("Total:           %d\n", s.Total)
	fmt.Printf("Estimated Time:  %d mins\n", s.EstimatedTime)
}

// ```

// ---

// ## Expected Output
// ```
// === Scenario 1: Normal Order ===
// Food Total:      298
// Drink Total:     90
// Dessert Total:   79
// Subtotal:        467
// Delivery Fee:    30
// Small Order Fee: 0
// Discount:        -29
// Service Fee:     23
// Total:           491
// Estimated Time:  25 mins

// === Scenario 2: Small Order + Express ===
// Food Total:      59
// Drink Total:     25
// Dessert Total:   0
// Subtotal:        84
// Delivery Fee:    0
// Small Order Fee: 20
// Discount:        0
// Service Fee:     5
// Total:           109
// Estimated Time:  15 mins

// === Scenario 3: Drinks Only ===
// Food Total:      0
// Drink Total:     255
// Dessert Total:   0
// Subtotal:        255
// Delivery Fee:    20
// Small Order Fee: 0
// Discount:        -51
// Service Fee:     11
// Total:           235
// Estimated Time:  20 mins
