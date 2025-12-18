package calculator

import (
	"math"
	"parkingapi/models"
)

var VehicleTypeHourlyRate = map[models.VehicleType]int{
	models.Motorcycle: 15,
	models.Car:        30,
	models.SUV:        40,
	models.Truck:      50,
}

type DiscountType string

var (
	EARLY_BIRD  DiscountType = "EARLY_BIRD"
	PERCENT_OFF DiscountType = "PERCENT_OFF"
	FLAT_OFF    DiscountType = "FLAT_OFF"
)

var ValidDiscountCode = map[string]bool{
	"EARLY_BIRD":    true,
	"MEMBER":        true,
	"MALL_PURCHASE": true,
}

var DiscountTypeCode = map[string]DiscountType{
	"EARLY_BIRD":    EARLY_BIRD,
	"MEMBER":        PERCENT_OFF,
	"MALL_PURCHASE": FLAT_OFF,
}

func GetHourlyRate(vehicleType models.VehicleType) int {
	return VehicleTypeHourlyRate[vehicleType]
}

func CalculateParkingFee(vehicleType models.VehicleType, minutes int) int {
	rate, exists := VehicleTypeHourlyRate[vehicleType]
	if !exists {
		return 0
	}

	hours := int(math.Ceil(float64(minutes) / 60))

	return hours * rate
}

func CalculateDiscount(code string, fee int) int {
	if _, exists := ValidDiscountCode[code]; !exists {
		return 0
	}

	discountType, exists := DiscountTypeCode[code]
	if !exists {
		return 0
	}

	switch discountType {
	case EARLY_BIRD:
		return fee
	case PERCENT_OFF:
		return int(float64(fee) * 0.2)
	case FLAT_OFF:
		return 50
	default:
		return 0
	}
}

func CalculateLostTicketFee(parkingFee int) int {
	if parkingFee < 0 {
		parkingFee = 0
	}

	return parkingFee + 100
}
