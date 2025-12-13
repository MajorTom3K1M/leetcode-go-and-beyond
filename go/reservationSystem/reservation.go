package reservation

import (
	"fmt"
	"slices"
	"time"
)

type Room struct {
	ID       string
	Name     string
	Capacity int
}

type Reservation struct {
	ID        string
	RoomID    string
	UserID    string
	Title     string
	StartTime time.Time
	EndTime   time.Time
}

type ReservationSystem struct {
	rooms          map[string]*Room
	reservations   map[string]*Reservation
	reservationSeq int
}

func NewReservationSystem() *ReservationSystem {
	return &ReservationSystem{
		rooms:          make(map[string]*Room),
		reservations:   make(map[string]*Reservation),
		reservationSeq: 0,
	}
}

// AddRoom adds a new room to the system.
// Returns false if room ID already exists or capacity <= 0.
func (rs *ReservationSystem) AddRoom(id, name string, capacity int) bool {
	_, ok := rs.rooms[id]
	if ok || capacity <= 0 {
		return false
	}
	rs.rooms[id] = &Room{
		ID:       id,
		Name:     name,
		Capacity: capacity,
	}
	return true
}

// GetRoom returns a room by ID, or nil if not found.
func (rs *ReservationSystem) GetRoom(id string) *Room {
	room, ok := rs.rooms[id]
	if ok {
		return room
	}
	return nil
}

// CreateReservation creates a new reservation.
// Returns reservation ID (format: "RES-{sequential number}") and true if successful.
// Returns empty string and false if:
//   - Room doesn't exist
//   - EndTime is not after StartTime
//   - Time slot conflicts with existing reservation for the same room
//
// Two reservations conflict if their time ranges overlap.
// Edge case: back-to-back is allowed (one ends exactly when another starts).
func (rs *ReservationSystem) CreateReservation(roomID, userID, title string, startTime, endTime time.Time) (string, bool) {
	_, ok := rs.rooms[roomID]
	if !ok {
		return "", false
	}

	if endTime.Sub(startTime) <= 0 {
		return "", false
	}

	result := []*Reservation{}
	for _, reserve := range rs.reservations {
		if reserve.RoomID == roomID {
			result = append(result, reserve)
		}
	}

	for _, reservation := range result {
		if startTime.Before(reservation.EndTime) && reservation.StartTime.Before(endTime) {
			return "", false
		}
	}

	rs.reservationSeq++
	reservationID := fmt.Sprintf("RES-%d", rs.reservationSeq)
	newReservation := &Reservation{
		ID:        reservationID,
		RoomID:    roomID,
		UserID:    userID,
		Title:     title,
		StartTime: startTime,
		EndTime:   endTime,
	}
	rs.reservations[reservationID] = newReservation

	return reservationID, true
}

// CancelReservation cancels a reservation by ID.
// Returns true if found and cancelled, false if not found.
func (rs *ReservationSystem) CancelReservation(id string) bool {
	if _, ok := rs.reservations[id]; !ok {
		return false
	}
	delete(rs.reservations, id)
	return true
}

// GetReservation returns a reservation by ID, or nil if not found.
func (rs *ReservationSystem) GetReservation(id string) *Reservation {
	if reservation, ok := rs.reservations[id]; ok {
		return reservation
	}
	return nil
}

// GetRoomReservations returns all reservations for a room on a given date.
// A reservation is "on" a date if any part of it falls on that date.
// Returns empty slice if room doesn't exist or no reservations found.
// Results should be sorted by StartTime ascending.
func (rs *ReservationSystem) GetRoomReservations(roomID string, date time.Time) []*Reservation {
	result := []*Reservation{}

	year, month, day := date.Date()
	dayStart := time.Date(year, month, day, 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.AddDate(0, 0, 1)

	for _, reservation := range rs.reservations {
		isInRange := reservation.StartTime.Before(dayEnd) && reservation.EndTime.After(dayStart)
		if reservation.RoomID == roomID && isInRange {
			result = append(result, reservation)
		}
	}

	slices.SortFunc(result, func(a, b *Reservation) int {
		return a.StartTime.Compare(b.StartTime)
	})

	return result
}

// GetUserReservations returns all reservations for a user.
// Returns empty slice if no reservations found.
func (rs *ReservationSystem) GetUserReservations(userID string) []*Reservation {
	result := []*Reservation{}
	for _, reservation := range rs.reservations {
		if reservation.UserID == userID {
			result = append(result, reservation)
		}
	}
	return result
}

// GetAvailableRooms returns all rooms that are available for the entire given time slot
// and have capacity >= minCapacity.
// Returns empty slice if none found.
func (rs *ReservationSystem) GetAvailableRooms(startTime, endTime time.Time, minCapacity int) []*Room {
	result := []*Room{}
	for _, room := range rs.rooms {
		if room.Capacity < minCapacity {
			continue
		}
		isAvailable := true
		for _, reservation := range rs.reservations {
			if reservation.RoomID == room.ID {
				if startTime.Before(reservation.EndTime) && reservation.StartTime.Before(endTime) {
					isAvailable = false
					break
				}
			}
		}
		if isAvailable {
			result = append(result, room)
		}
	}
	return result
}

// UpdateReservation updates an existing reservation's time slot.
// Returns true if successful, false if:
//   - Reservation doesn't exist
//   - New time slot is invalid (end not after start)
//   - New time slot conflicts with other reservations (excluding itself)
func (rs *ReservationSystem) UpdateReservation(reservationID string, newStart, newEnd time.Time) bool {
	reservation, ok := rs.reservations[reservationID]
	if !ok {
		return false
	}
	if newEnd.Sub(newStart) <= 0 {
		return false
	}
	for _, otherReservation := range rs.reservations {
		if otherReservation.RoomID == reservation.RoomID && otherReservation.ID != reservationID {
			if newStart.Before(otherReservation.EndTime) && otherReservation.StartTime.Before(newEnd) {
				return false
			}
		}
	}
	reservation.StartTime = newStart
	reservation.EndTime = newEnd
	return true
}
