package spot

import (
	"fmt"
	m "parkingapi/models"
)

var CanFitMap = map[m.SpotSize]map[m.VehicleType]bool{
	m.Small: {
		m.Motorcycle: true,
	},
	m.Medium: {
		m.Motorcycle: true,
		m.Car:        true,
	},
	m.Large: {
		m.Motorcycle: true,
		m.Car:        true,
		m.SUV:        true,
	},
	m.XLarge: {
		m.Motorcycle: true,
		m.Car:        true,
		m.SUV:        true,
		m.Truck:      true,
	},
}

var SpotSizePriority = []m.SpotSize{
	m.Small,
	m.Medium,
	m.Large,
	m.XLarge,
}

func CanFit(vehicleType m.VehicleType, spotSize m.SpotSize) bool {
	if fitMap, exists := CanFitMap[spotSize]; exists {
		if _, exists := fitMap[vehicleType]; exists {
			return true
		}
	}

	return false
}

func GetSpotPriority(vehicleType m.VehicleType) []m.SpotSize {
	out := make([]m.SpotSize, 0, len(SpotSizePriority))

	for _, size := range SpotSizePriority {
		fitMap, ok := CanFitMap[size]
		if !ok {
			continue
		}
		if fitMap[vehicleType] { // true means can fit
			out = append(out, size)
		}
	}

	return out
}

func GenerateSpotID(floor int, size m.SpotSize, number int) string {
	spotID := fmt.Sprintf("F%d-%s%02d", floor, size, number)
	return spotID
}
