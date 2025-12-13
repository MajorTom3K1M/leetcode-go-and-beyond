package reservation

import (
	"testing"
	"time"
)

func makeTime(hour, minute int) time.Time {
	return time.Date(2024, 1, 15, hour, minute, 0, 0, time.UTC)
}

func makeDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func TestAddAndGetRoom(t *testing.T) {
	rs := NewReservationSystem()

	// Add room successfully
	if !rs.AddRoom("room-1", "Conference A", 10) {
		t.Error("should add room successfully")
	}

	// Get room
	room := rs.GetRoom("room-1")
	if room == nil {
		t.Fatal("room should not be nil")
	}
	if room.Name != "Conference A" || room.Capacity != 10 {
		t.Error("room data mismatch")
	}

	// Duplicate ID should fail
	if rs.AddRoom("room-1", "Different Name", 5) {
		t.Error("duplicate room ID should fail")
	}

	// Invalid capacity should fail
	if rs.AddRoom("room-2", "Small Room", 0) {
		t.Error("zero capacity should fail")
	}
	if rs.AddRoom("room-3", "Negative Room", -5) {
		t.Error("negative capacity should fail")
	}

	// Non-existent room
	if rs.GetRoom("room-999") != nil {
		t.Error("non-existent room should return nil")
	}
}

func TestCreateReservation(t *testing.T) {
	rs := NewReservationSystem()
	rs.AddRoom("room-1", "Conference A", 10)

	id, ok := rs.CreateReservation("room-1", "user-1", "Team Meeting",
		makeTime(9, 0), makeTime(10, 0))
	if !ok || id != "RES-1" {
		t.Errorf("expected RES-1 and true, got %s and %v", id, ok)
	}

	res := rs.GetReservation(id)
	if res == nil {
		t.Fatal("reservation should exist")
	}
	if res.RoomID != "room-1" || res.UserID != "user-1" || res.Title != "Team Meeting" {
		t.Error("reservation data mismatch")
	}

	id2, ok := rs.CreateReservation("room-1", "user-2", "Another Meeting",
		makeTime(11, 0), makeTime(12, 0))
	if !ok || id2 != "RES-2" {
		t.Errorf("expected RES-2 and true, got %s and %v", id2, ok)
	}

	_, ok = rs.CreateReservation("room-999", "user-1", "Test", makeTime(9, 0), makeTime(10, 0))
	if ok {
		t.Error("non-existent room should fail")
	}

	_, ok = rs.CreateReservation("room-1", "user-1", "Test", makeTime(10, 0), makeTime(9, 0))
	if ok {
		t.Error("end before start should fail")
	}

	_, ok = rs.CreateReservation("room-1", "user-1", "Test", makeTime(10, 0), makeTime(10, 0))
	if ok {
		t.Error("same start and end should fail")
	}
}

func TestReservationConflicts(t *testing.T) {
	rs := NewReservationSystem()
	rs.AddRoom("room-1", "Conference A", 10)

	// Create base reservation: 10:00 - 12:00
	rs.CreateReservation("room-1", "user-1", "Base Meeting",
		makeTime(10, 0), makeTime(12, 0))

	// Exact same time should conflict
	_, ok := rs.CreateReservation("room-1", "user-2", "Conflict",
		makeTime(10, 0), makeTime(12, 0))
	if ok {
		t.Error("exact overlap should conflict")
	}

	// Starts during existing should conflict
	_, ok = rs.CreateReservation("room-1", "user-2", "Conflict",
		makeTime(11, 0), makeTime(13, 0))
	if ok {
		t.Error("starting during existing should conflict")
	}

	// Ends during existing should conflict
	_, ok = rs.CreateReservation("room-1", "user-2", "Conflict",
		makeTime(9, 0), makeTime(11, 0))
	if ok {
		t.Error("ending during existing should conflict")
	}

	// Completely contains existing should conflict
	_, ok = rs.CreateReservation("room-1", "user-2", "Conflict",
		makeTime(9, 0), makeTime(13, 0))
	if ok {
		t.Error("containing existing should conflict")
	}

	// Inside existing should conflict
	_, ok = rs.CreateReservation("room-1", "user-2", "Conflict",
		makeTime(10, 30), makeTime(11, 30))
	if ok {
		t.Error("inside existing should conflict")
	}

	// Back-to-back should NOT conflict (ends exactly when other starts)
	_, ok = rs.CreateReservation("room-1", "user-2", "After",
		makeTime(12, 0), makeTime(13, 0))
	if !ok {
		t.Error("back-to-back (after) should not conflict")
	}

	_, ok = rs.CreateReservation("room-1", "user-2", "Before",
		makeTime(8, 0), makeTime(10, 0))
	if !ok {
		t.Error("back-to-back (before) should not conflict")
	}

	// Different room should NOT conflict
	rs.AddRoom("room-2", "Conference B", 8)
	_, ok = rs.CreateReservation("room-2", "user-2", "Different Room",
		makeTime(10, 0), makeTime(12, 0))
	if !ok {
		t.Error("different room should not conflict")
	}
}

func TestCancelReservation(t *testing.T) {
	rs := NewReservationSystem()
	rs.AddRoom("room-1", "Conference A", 10)

	id, _ := rs.CreateReservation("room-1", "user-1", "Meeting",
		makeTime(10, 0), makeTime(11, 0))

	// Cancel existing
	if !rs.CancelReservation(id) {
		t.Error("should cancel existing reservation")
	}

	// Should be gone
	if rs.GetReservation(id) != nil {
		t.Error("cancelled reservation should not exist")
	}

	// Cancel non-existent
	if rs.CancelReservation("RES-999") {
		t.Error("cancelling non-existent should return false")
	}

	// After cancellation, time slot should be available again
	_, ok := rs.CreateReservation("room-1", "user-2", "New Meeting",
		makeTime(10, 0), makeTime(11, 0))
	if !ok {
		t.Error("should be able to book cancelled slot")
	}
}

func TestGetRoomReservations(t *testing.T) {
	rs := NewReservationSystem()
	rs.AddRoom("room-1", "Conference A", 10)

	rs.CreateReservation("room-1", "user-1", "Morning",
		makeTime(9, 0), makeTime(10, 0))
	rs.CreateReservation("room-1", "user-1", "Afternoon",
		makeTime(14, 0), makeTime(15, 0))
	rs.CreateReservation("room-1", "user-1", "Midday",
		makeTime(11, 0), makeTime(12, 0))

	rs.CreateReservation("room-1", "user-1", "Tomorrow",
		time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 16, 11, 0, 0, 0, time.UTC))

	reservations := rs.GetRoomReservations("room-1", makeDate(2024, 1, 15))
	if len(reservations) != 3 {
		t.Errorf("expected 3 reservations, got %d", len(reservations))
	}

	if reservations[0].Title != "Morning" ||
		reservations[1].Title != "Midday" ||
		reservations[2].Title != "Afternoon" {
		t.Error("reservations should be sorted by start time")
	}

	empty := rs.GetRoomReservations("room-999", makeDate(2024, 1, 15))
	if len(empty) != 0 {
		t.Error("non-existent room should return empty slice")
	}

	empty = rs.GetRoomReservations("room-1", makeDate(2024, 1, 20))
	if len(empty) != 0 {
		t.Error("day with no reservations should return empty slice")
	}
}

func TestGetUserReservations(t *testing.T) {
	rs := NewReservationSystem()
	rs.AddRoom("room-1", "Conference A", 10)
	rs.AddRoom("room-2", "Conference B", 8)

	rs.CreateReservation("room-1", "user-1", "Meeting 1", makeTime(9, 0), makeTime(10, 0))
	rs.CreateReservation("room-2", "user-1", "Meeting 2", makeTime(11, 0), makeTime(12, 0))
	rs.CreateReservation("room-1", "user-2", "Other User", makeTime(13, 0), makeTime(14, 0))

	user1Res := rs.GetUserReservations("user-1")
	if len(user1Res) != 2 {
		t.Errorf("expected 2 reservations for user-1, got %d", len(user1Res))
	}

	user2Res := rs.GetUserReservations("user-2")
	if len(user2Res) != 1 {
		t.Errorf("expected 1 reservation for user-2, got %d", len(user2Res))
	}

	unknownRes := rs.GetUserReservations("user-999")
	if len(unknownRes) != 0 {
		t.Error("unknown user should have empty slice")
	}
}

func TestGetAvailableRooms(t *testing.T) {
	rs := NewReservationSystem()
	rs.AddRoom("small", "Small Room", 4)
	rs.AddRoom("medium", "Medium Room", 8)
	rs.AddRoom("large", "Large Room", 20)

	rs.CreateReservation("medium", "user-1", "Meeting",
		makeTime(10, 0), makeTime(12, 0))

	available := rs.GetAvailableRooms(makeTime(10, 0), makeTime(11, 0), 5)
	if len(available) != 1 {
		t.Errorf("expected 1 available room, got %d", len(available))
	}
	if available[0].ID != "large" {
		t.Errorf("expected large room, got %s", available[0].ID)
	}

	available = rs.GetAvailableRooms(makeTime(13, 0), makeTime(14, 0), 1)
	if len(available) != 3 {
		t.Errorf("expected 3 available rooms, got %d", len(available))
	}

	available = rs.GetAvailableRooms(makeTime(13, 0), makeTime(14, 0), 100)
	if len(available) != 0 {
		t.Error("no room should have capacity 100")
	}
}

func TestUpdateReservation(t *testing.T) {
	rs := NewReservationSystem()
	rs.AddRoom("room-1", "Conference A", 10)

	id1, _ := rs.CreateReservation("room-1", "user-1", "First",
		makeTime(9, 0), makeTime(10, 0))
	rs.CreateReservation("room-1", "user-2", "Second",
		makeTime(11, 0), makeTime(12, 0))

	if !rs.UpdateReservation(id1, makeTime(13, 0), makeTime(14, 0)) {
		t.Error("valid update should succeed")
	}
	res := rs.GetReservation(id1)
	if res.StartTime != makeTime(13, 0) || res.EndTime != makeTime(14, 0) {
		t.Error("times should be updated")
	}

	if rs.UpdateReservation(id1, makeTime(11, 0), makeTime(12, 0)) {
		t.Error("conflicting update should fail")
	}

	if rs.UpdateReservation("RES-999", makeTime(15, 0), makeTime(16, 0)) {
		t.Error("updating non-existent should fail")
	}

	if rs.UpdateReservation(id1, makeTime(15, 0), makeTime(14, 0)) {
		t.Error("end before start should fail")
	}

	if !rs.UpdateReservation(id1, makeTime(13, 0), makeTime(14, 30)) {
		t.Error("extending own reservation should succeed if no conflict")
	}
}
