package models

import "time"

type VehicleType string

const (
	Motorcycle VehicleType = "motorcycle"
	Car        VehicleType = "car"
	SUV        VehicleType = "suv"
	Truck      VehicleType = "truck"
)

type SpotSize string

const (
	Small  SpotSize = "small"  // motorcycle only
	Medium SpotSize = "medium" // motorcycle, car
	Large  SpotSize = "large"  // motorcycle, car, suv
	XLarge SpotSize = "xlarge" // all vehicles
)

type Vehicle struct {
	PlateNumber string
	VehicleType VehicleType
}

type ParkingSpot struct {
	SpotID     string
	Size       SpotSize
	Floor      int
	IsOccupied bool
	Vehicle    *Vehicle
}

type ParkingTicket struct {
	TicketID  string
	Vehicle   Vehicle
	SpotID    string
	Floor     int
	EntryTime time.Time
	ExitTime  *time.Time
	Duration  int // minutes
	Fee       int
	IsPaid    bool
}

type ParkingReceipt struct {
	TicketID    string
	PlateNumber string
	VehicleType VehicleType
	SpotID      string
	Floor       int
	EntryTime   time.Time
	ExitTime    time.Time
	Duration    int
	HourlyRate  int
	ParkingFee  int
	Discount    int
	TotalFee    int
}

type ParkingStatus struct {
	TotalSpots      int
	AvailableSpots  int
	OccupiedSpots   int
	SmallAvailable  int
	MediumAvailable int
	LargeAvailable  int
	XLargeAvailable int
}
