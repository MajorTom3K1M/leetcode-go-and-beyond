package main

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

var (
	ErrMovieAlreadyExists    = errors.New("movie already exists")
	ErrScreenAlreadyExists   = errors.New("screen already exists")
	ErrShowtimeAlreadyExists = errors.New("showtime already exists")
	ErrShowtimeNotFound      = errors.New("showtime not found")
	ErrScreenNotFound        = errors.New("screen not found")
	ErrSeatIsNotAvailable    = errors.New("seat not available")
	ErrBookingNotFound       = errors.New("booking not found")
	ErrMovieNotFound         = errors.New("movie not found")
)

type Movie struct {
	ID          string
	Title       string
	DurationMin int
}

type Screen struct {
	ID          string
	Rows        int
	SeatsPerRow int
}

type Showtime struct {
	ID        string
	MovieID   string
	ScreenID  string
	StartTime time.Time
	Booked    map[string]bool
}

type Booking struct {
	ID           string
	ShowtimeID   string
	CustomerName string
	Seats        []string
	CreatedAt    time.Time
}

type Cinema struct {
	Movies         map[string]*Movie    // movieID -> Movie
	Screens        map[string]*Screen   // screenID -> Screen
	Showtimes      map[string]*Showtime // showtimeID -> Showtime
	Bookings       map[string]*Booking  // bookingID -> Booking
	bookingCounter int
}

func NewCinema() *Cinema {
	return &Cinema{
		Movies:    make(map[string]*Movie),
		Screens:   make(map[string]*Screen),
		Showtimes: make(map[string]*Showtime),
		Bookings:  make(map[string]*Booking),
	}
}

func (c *Cinema) AddMovie(id, title string, durationMin int) error {
	_, ok := c.Movies[id]
	if ok {
		return ErrMovieAlreadyExists
	}

	c.Movies[id] = &Movie{
		ID:          id,
		Title:       title,
		DurationMin: durationMin,
	}

	return nil
}

func (c *Cinema) AddScreen(id string, rows, seatsPerRow int) error {
	_, ok := c.Screens[id]
	if ok {
		return ErrScreenAlreadyExists
	}

	c.Screens[id] = &Screen{
		ID:          id,
		Rows:        rows,
		SeatsPerRow: seatsPerRow,
	}

	return nil
}

func (c *Cinema) AddShowtime(id, movieID, screenID string, startTime time.Time) error {
	c.CleanUpShowtime()

	if _, exists := c.Showtimes[id]; exists {
		return ErrShowtimeAlreadyExists
	}

	if _, exists := c.Movies[movieID]; !exists {
		return ErrMovieNotFound
	}

	// Create Available Seat
	screen, ok := c.Screens[screenID]
	if !ok {
		return ErrScreenNotFound
	}

	booked := map[string]bool{}
	for row := 0; row < screen.Rows; row++ {
		rowLetter := string(rune('A' + row))
		for seat := 1; seat <= screen.SeatsPerRow; seat++ {
			seatID := fmt.Sprintf("%s%d", rowLetter, seat)
			booked[seatID] = false
		}
	}

	c.Showtimes[id] = &Showtime{
		ID:        id,
		MovieID:   movieID,
		ScreenID:  screenID,
		StartTime: startTime,
		Booked:    booked,
	}

	return nil
}

func (c *Cinema) GetAvailableSeats(showtimeID string) ([]string, error) {
	showtime, ok := c.Showtimes[showtimeID]
	if !ok {
		return nil, ErrShowtimeNotFound
	}

	availableSeats := []string{}
	for seat, booked := range showtime.Booked {
		if !booked {
			availableSeats = append(availableSeats, seat)
		}
	}

	sort.Strings(availableSeats)
	return availableSeats, nil
}

func (c *Cinema) BookSeats(showtimeID string, seats []string, customerName string) (string, error) {
	showtime, ok := c.Showtimes[showtimeID]
	if !ok {
		return "", ErrShowtimeNotFound
	}

	for _, seat := range seats {
		isBooked, ok := showtime.Booked[seat]
		if !ok || isBooked {
			return "", ErrSeatIsNotAvailable
		}
	}

	for _, seat := range seats {
		showtime.Booked[seat] = true
	}

	c.bookingCounter++
	bookingID := fmt.Sprintf("BK%d", c.bookingCounter)

	// if duplicate booking just update it !!
	c.Bookings[bookingID] = &Booking{
		ID:           bookingID,
		ShowtimeID:   showtimeID,
		CustomerName: customerName,
		Seats:        seats,
		CreatedAt:    time.Now(),
	}

	return customerName, nil
}

func (c *Cinema) CleanUpShowtime() {
	now := time.Now()
	endedShow := []string{}
	for _, showtime := range c.Showtimes {
		movie, ok := c.Movies[showtime.MovieID]
		if !ok || showtime.StartTime.Add(time.Duration(movie.DurationMin)*time.Minute).Before(now) {
			endedShow = append(endedShow, showtime.ID)
		}
	}

	for _, ended := range endedShow {
		delete(c.Showtimes, ended)
	}
}

func (c *Cinema) CancelBooking(bookingID string) error {
	booking, ok := c.Bookings[bookingID]
	if !ok {
		return ErrBookingNotFound
	}

	showtime, ok := c.Showtimes[booking.ShowtimeID]
	if ok {
		for _, seat := range booking.Seats {
			showtime.Booked[seat] = false
		}
	}

	delete(c.Bookings, bookingID)
	return nil
}

func (c *Cinema) GetBooking(bookingID string) (*Booking, error) {
	booking, ok := c.Bookings[bookingID]
	if !ok {
		return nil, ErrBookingNotFound
	}

	return booking, nil
}
