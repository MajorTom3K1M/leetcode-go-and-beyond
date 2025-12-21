package main

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrAlreadyReserved      = errors.New("room already occupied")
	ErrReserveAlreadyExists = errors.New("reservation already exists")
)

type Room struct {
	ID            string
	RoomType      string
	PricePerNight float64
}

type Reservation struct {
	ID        string
	RoomID    string
	GuestName string
	CheckIn   time.Time
	CheckOut  time.Time
}

type Hotel struct {
	Rooms            map[string]*Room
	Reservations     map[string]*Reservation
	reservationCount int
}

func NewHotel() *Hotel {
	return &Hotel{
		Rooms:        make(map[string]*Room),
		Reservations: make(map[string]*Reservation),
	}
}

func (h *Hotel) AddRoom(id string, roomType string, pricePerNight float64) {
	if _, ok := h.Rooms[id]; ok {
		return
	}

	h.Rooms[id] = &Room{
		ID:            id,
		RoomType:      roomType,
		PricePerNight: pricePerNight,
	}
}

func (h *Hotel) GetAvailableRooms(checkIn, checkOut time.Time, roomType string) []string {
	availableRooms := []string{}
	for _, room := range h.Rooms {
		if room.RoomType != roomType {
			continue
		}

		if h.isRoomAvailable(room.ID, checkIn, checkOut) {
			availableRooms = append(availableRooms, room.ID)
		}
	}

	return availableRooms
}

func (h *Hotel) isRoomAvailable(roomID string, checkIn, checkOut time.Time) bool {
	y1, m1, d1 := checkIn.Date()
	y2, m2, d2 := checkOut.Date()
	newCheckIn := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	newCheckOut := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)

	for _, res := range h.Reservations {
		if res.RoomID != roomID {
			continue
		}

		y3, m3, d3 := res.CheckIn.Date()
		y4, m4, d4 := res.CheckOut.Date()
		reserveCheckIn := time.Date(y3, m3, d3, 0, 0, 0, 0, time.UTC)
		reserveCheckOut := time.Date(y4, m4, d4, 0, 0, 0, 0, time.UTC)

		if IsOverlap(newCheckIn, newCheckOut, reserveCheckIn, reserveCheckOut) {
			return false
		}
	}
	return true
}

func IsOverlap(aCheckIn, aCheckOut, bCheckIn, bCheckout time.Time) bool {
	if aCheckIn.Before(bCheckout) && aCheckOut.After(bCheckIn) {
		return true
	}
	return false
}

func (h *Hotel) MakeReservation(roomID, guestName string, checkIn, checkOut time.Time) (reservationID string, err error) {
	for _, revertation := range h.Reservations {
		if revertation.RoomID == roomID {
			y1, m1, d1 := checkIn.Date()
			y2, m2, d2 := checkOut.Date()
			y3, m3, d3 := revertation.CheckIn.Date()
			y4, m4, d4 := revertation.CheckOut.Date()
			newCheckIn := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
			newCheckOut := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)
			reserveCheckIn := time.Date(y3, m3, d3, 0, 0, 0, 0, time.UTC)
			reserveCheckOut := time.Date(y4, m4, d4, 0, 0, 0, 0, time.UTC)
			if IsOverlap(newCheckIn, newCheckOut, reserveCheckIn, reserveCheckOut) {
				return "", ErrAlreadyReserved
			}
		}
	}

	h.reservationCount++
	resID := fmt.Sprintf("%d", h.reservationCount)
	if _, ok := h.Reservations[resID]; ok {
		return "", ErrReserveAlreadyExists
	}

	h.Reservations[resID] = &Reservation{
		ID:        resID,
		RoomID:    roomID,
		GuestName: guestName,
		CheckIn:   checkIn,
		CheckOut:  checkOut,
	}

	return resID, nil
}

func (h *Hotel) CancelReservation(reservationID string) {
	if _, ok := h.Reservations[reservationID]; !ok {
		return
	}

	delete(h.Reservations, reservationID)
}

func (h *Hotel) GetReservation(reservationID string) *Reservation {
	if _, ok := h.Reservations[reservationID]; !ok {
		return nil
	}

	return h.Reservations[reservationID]
}

func (h *Hotel) GetGuestReservations(guestName string) []*Reservation {
	guestReservations := []*Reservation{}
	for _, reservation := range h.Reservations {
		if reservation.GuestName == guestName {
			guestReservations = append(guestReservations, reservation)
		}
	}

	return guestReservations
}

func (h *Hotel) CalculateTotalPrice(reservationID string) float64 {
	reservation := h.GetReservation(reservationID)

	duration := reservation.CheckOut.Sub(reservation.CheckIn)
	days := int((duration + 24*time.Hour - 1) / (24 * time.Hour))

	room, ok := h.Rooms[reservation.RoomID]
	if !ok {
		return 0
	}

	return room.PricePerNight * float64(days)
}

func main() {
	hotel := NewHotel()

	hotel.AddRoom("101", "single", 100.0)
	hotel.AddRoom("102", "single", 100.0)
	hotel.AddRoom("201", "double", 150.0)
	hotel.AddRoom("301", "suite", 300.0)

	checkIn := time.Date(2024, 3, 15, 14, 0, 0, 0, time.Local)
	checkOut := time.Date(2024, 3, 18, 11, 0, 0, 0, time.Local) // 3 nights

	// Find available singles
	rooms := hotel.GetAvailableRooms(checkIn, checkOut, "single") // ["101", "102"]
	fmt.Println("Rooms : ", rooms)

	// Book room 101
	resID, _ := hotel.MakeReservation("101", "John Doe", checkIn, checkOut)

	// Room 101 no longer available for overlapping dates
	rooms = hotel.GetAvailableRooms(checkIn, checkOut, "single") // ["102"]
	fmt.Println("Rooms : ", rooms)

	// But available for different dates
	otherCheckIn := time.Date(2024, 3, 20, 14, 0, 0, 0, time.Local)
	otherCheckOut := time.Date(2024, 3, 22, 11, 0, 0, 0, time.Local)
	rooms = hotel.GetAvailableRooms(otherCheckIn, otherCheckOut, "single") // ["101", "102"]
	fmt.Println("Rooms : ", rooms)

	// Calculate price
	total := hotel.CalculateTotalPrice(resID) // 300.0 (3 nights × $100)
	fmt.Println("total : ", total)

	// Get guest history
	hotel.MakeReservation("201", "John Doe", otherCheckIn, otherCheckOut)
	reservations := hotel.GetGuestReservations("John Doe") // 2 reservations

	for _, reserve := range reservations {
		fmt.Printf("reservations : %+v\n", reserve)
	}
}
