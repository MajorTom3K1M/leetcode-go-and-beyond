package meetingroomscheduler

import "testing"

func TestBookSimpleNonOverlap(t *testing.T) {
	s := NewScheduler()

	if !s.Book(Meeting{Start: 9 * 60, End: 10 * 60, Title: "Standup"}) {
		t.Fatal("expected first booking to succeed")
	}
	if !s.Book(Meeting{Start: 10 * 60, End: 11 * 60, Title: "Planning"}) {
		t.Fatal("expected second booking to succeed")
	}

	ms := s.Meetings()
	if len(ms) != 2 {
		t.Fatalf("expected 2 meetings, got %d", len(ms))
	}

	if ms[0].Title != "Standup" || ms[1].Title != "Planning" {
		t.Fatalf("unexpected order: %+v", ms)
	}
}

func TestOverlapShouldFail(t *testing.T) {
	s := NewScheduler()

	s.Book(Meeting{Start: 9 * 60, End: 10 * 60, Title: "Standup"})

	if s.Book(Meeting{Start: 9*60 + 30, End: 9*60 + 45, Title: "Conflict"}) {
		t.Fatal("expected overlapping meeting to be rejected")
	}
}

func TestTouchingIsAllowed(t *testing.T) {
	s := NewScheduler()

	if !s.Book(Meeting{Start: 9 * 60, End: 10 * 60, Title: "A"}) {
		t.Fatal("expected first booking to succeed")
	}
	if !s.Book(Meeting{Start: 10 * 60, End: 11 * 60, Title: "B"}) {
		t.Fatal("expected non-overlapping touching booking to succeed")
	}
}

func TestInvalidMeeting(t *testing.T) {
	s := NewScheduler()
	if s.Book(Meeting{Start: 10 * 60, End: 9 * 60, Title: "Invalid"}) {
		t.Fatal("expected invalid meeting to be rejected")
	}
}
