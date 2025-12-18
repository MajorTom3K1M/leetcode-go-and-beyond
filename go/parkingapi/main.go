package main

import (
	"fmt"
	"parkingapi/models"
	"parkingapi/parking"
)

func main() {
	// Initialize parking lot
	// Floor 1: 2 small, 3 medium, 2 large, 1 xlarge
	// Floor 2: 1 small, 2 medium, 3 large, 2 xlarge
	lot := parking.NewParkingLot("Mall Parking")

	lot.AddSpots(1, "small", 2)
	lot.AddSpots(1, "medium", 3)
	lot.AddSpots(1, "large", 2)
	lot.AddSpots(1, "xlarge", 1)

	lot.AddSpots(2, "small", 1)
	lot.AddSpots(2, "medium", 2)
	lot.AddSpots(2, "large", 3)
	lot.AddSpots(2, "xlarge", 2)

	fmt.Println("=== Initial Status ===")
	printStatus(lot.GetStatus())

	// Scenario 1: Cars entering
	fmt.Println("\n=== Vehicles Entering ===")

	ticket1, _ := lot.Entry("ABC-1234", "car")
	fmt.Printf("Car ABC-1234 parked at %s (Floor %d)\n", ticket1.SpotID, ticket1.Floor)

	ticket2, _ := lot.Entry("XYZ-5678", "motorcycle")
	fmt.Printf("Motorcycle XYZ-5678 parked at %s (Floor %d)\n", ticket2.SpotID, ticket2.Floor)

	ticket3, _ := lot.Entry("DEF-9999", "suv")
	fmt.Printf("SUV DEF-9999 parked at %s (Floor %d)\n", ticket3.SpotID, ticket3.Floor)

	ticket4, err := lot.Entry("GHI-1111", "truck")
	if err != nil {
		fmt.Printf("Truck GHI-1111 could not be parked: %v\n", err)
	} else {
		fmt.Printf("Truck GHI-1111 parked at %s (Floor %d)\n", ticket4.SpotID, ticket4.Floor)
	}

	fmt.Println("\n=== Status After Entry ===")
	printStatus(lot.GetStatus())

	// Scenario 2: Simulate time passing and exit
	fmt.Println("\n=== Vehicles Exiting ===")

	// Simulate 2 hours 30 minutes parking for car
	receipt1, _ := lot.Exit(ticket1.TicketID, 150) // 150 minutes
	printReceipt(receipt1)

	// Simulate 45 minutes parking for motorcycle
	receipt2, _ := lot.Exit(ticket2.TicketID, 45)
	printReceipt(receipt2)

	// Scenario 3: With mall validation (discount)
	fmt.Println("\n=== Exit with Mall Discount ===")

	// SUV with 3 hours parking + mall purchase discount
	receipt3, _ := lot.ExitWithDiscount(ticket3.TicketID, 180, "MALL_PURCHASE")
	printReceipt(receipt3)

	// Scenario 4: Lost ticket
	fmt.Println("\n=== Lost Ticket ===")
	receipt4, _ := lot.ExitLostTicket("GHI-1111", 240) // 4 hours, lost ticket fee
	printReceipt(receipt4)

	fmt.Println("\n=== Final Status ===")
	printStatus(lot.GetStatus())
}

func printStatus(s *models.ParkingStatus) {
	fmt.Printf("Total: %d | Available: %d | Occupied: %d\n",
		s.TotalSpots, s.AvailableSpots, s.OccupiedSpots)
	fmt.Printf("Small: %d | Medium: %d | Large: %d | XLarge: %d\n",
		s.SmallAvailable, s.MediumAvailable, s.LargeAvailable, s.XLargeAvailable)
}

func printReceipt(r *models.ParkingReceipt) {
	fmt.Println("-------------------")
	fmt.Printf("Ticket: %s\n", r.TicketID)
	fmt.Printf("Plate: %s (%s)\n", r.PlateNumber, r.VehicleType)
	fmt.Printf("Spot: %s (Floor %d)\n", r.SpotID, r.Floor)
	fmt.Printf("Duration: %d mins\n", r.Duration)
	fmt.Printf("Rate: %d baht/hr\n", r.HourlyRate)
	fmt.Printf("Parking Fee: %d baht\n", r.ParkingFee)
	fmt.Printf("Discount: -%d baht\n", r.Discount)
	fmt.Printf("Total: %d baht\n", r.TotalFee)
	fmt.Println("-------------------")
}

// ```

// ---

// ## Expected Output
// ```
// === Initial Status ===
// Total: 16 | Available: 16 | Occupied: 0
// Small: 3 | Medium: 5 | Large: 5 | XLarge: 3

// === Vehicles Entering ===
// Car ABC-1234 parked at F1-M01 (Floor 1)
// Motorcycle XYZ-5678 parked at F1-S01 (Floor 1)
// SUV DEF-9999 parked at F1-L01 (Floor 1)
// Truck GHI-1111 parked at F1-XL01 (Floor 1)

// === Status After Entry ===
// Total: 16 | Available: 12 | Occupied: 4
// Small: 2 | Medium: 4 | Large: 4 | XLarge: 2

// === Vehicles Exiting ===
// -------------------
// Ticket: T001
// Plate: ABC-1234 (car)
// Spot: F1-M01 (Floor 1)
// Duration: 150 mins
// Rate: 30 baht/hr
// Parking Fee: 90 baht
// Discount: -0 baht
// Total: 90 baht
// -------------------
// -------------------
// Ticket: T002
// Plate: XYZ-5678 (motorcycle)
// Spot: F1-S01 (Floor 1)
// Duration: 45 mins
// Rate: 15 baht/hr
// Parking Fee: 15 baht
// Discount: -0 baht
// Total: 15 baht
// -------------------

// === Exit with Mall Discount ===
// -------------------
// Ticket: T003
// Plate: DEF-9999 (suv)
// Spot: F1-L01 (Floor 1)
// Duration: 180 mins
// Rate: 40 baht/hr
// Parking Fee: 120 baht
// Discount: -50 baht
// Total: 70 baht
// -------------------

// === Lost Ticket ===
// -------------------
// Ticket: T004
// Plate: GHI-1111 (truck)
// Spot: F1-XL01 (Floor 1)
// Duration: 240 mins
// Rate: 50 baht/hr
// Parking Fee: 200 baht
// Discount: -0 baht
// Total: 300 baht
// -------------------

// === Final Status ===
// Total: 16 | Available: 16 | Occupied: 0
// Small: 3 | Medium: 5 | Large: 5 | XLarge: 3
// ```

// ---

// ## Business Rules

// ### Hourly Rates (calculator/calculator.go)

// | Vehicle Type | Rate per Hour |
// |--------------|---------------|
// | Motorcycle | 15 baht |
// | Car | 30 baht |
// | SUV | 40 baht |
// | Truck | 50 baht |

// ### Fee Calculation

// - First hour: full rate
// - After first hour: charge per hour (round up)
// - Example: 150 mins = 3 hours → 3 × rate

// ### Spot Assignment Rules (spot/spot.go)

// | Vehicle Type | Can Park In |
// |--------------|-------------|
// | Motorcycle | small, medium, large, xlarge |
// | Car | medium, large, xlarge |
// | SUV | large, xlarge |
// | Truck | xlarge only |

// **Priority:** Assign smallest available spot first.

// ### Spot ID Format
// ```
// F{floor}-{size}{number}
// Examples: F1-S01, F1-M01, F2-L03, F2-XL01
