package meetingroomscheduler

import "sort"

type Meeting struct {
	Start int // minutes from 0 to 24*60, e.g. 9:30 = 570
	End   int // strictly greater than Start
	Title string
}

type Scheduler struct {
	books []Meeting
}

func NewScheduler() *Scheduler {
	return &Scheduler{books: []Meeting{}}
}

// Book tries to add a meeting.
// Returns true if booked successfully (no conflict),
// or false if it overlaps any existing meeting.
func (s *Scheduler) Book(m Meeting) bool {
	if m.End <= m.Start {
		return false
	}

	if m.Start < 0 || m.End > 24*60 {
		return false
	}

	for _, booked := range s.books {
		if m.Start < booked.End && booked.Start < m.End {
			return false
		}
	}
	s.books = append(s.books, m)
	return true
}

// Meetings returns all meetings sorted by start time.
func (s *Scheduler) Meetings() []Meeting {
	sort.Slice(s.books, func(i, j int) bool {
		return s.books[i].Start < s.books[j].Start
	})
	return s.books
}
