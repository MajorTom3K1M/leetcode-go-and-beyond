package leaderboard

import "testing"

func TestEmptyLeaderboardTop(t *testing.T) {
	lb := NewLeaderboard()
	if got := lb.Top(3); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestSinglePlayer(t *testing.T) {
	lb := NewLeaderboard()
	lb.AddScore("alice", 50)

	if got := lb.Top(1); got != 50 {
		t.Fatalf("expected 50, got %d", got)
	}

	lb.Reset("alice")
	if got := lb.Top(1); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestMultiplePlayersTopK(t *testing.T) {
	lb := NewLeaderboard()

	lb.AddScore("alice", 50)
	lb.AddScore("bob", 30)
	lb.AddScore("charlie", 70)

	if got := lb.Top(2); got != 120 {
		t.Fatalf("expected 120, got %d", got)
	}

	if got := lb.Top(5); got != 150 {
		t.Fatalf("expected 150, got %d", got)
	}
}

func TestAddScoreAccumulates(t *testing.T) {
	lb := NewLeaderboard()

	lb.AddScore("alice", 10)
	lb.AddScore("alice", 20)

	if got := lb.Top(1); got != 30 {
		t.Fatalf("expected 30, got %d", got)
	}
}

func TestResetPlayer(t *testing.T) {
	lb := NewLeaderboard()

	lb.AddScore("alice", 40)
	lb.AddScore("bob", 60)

	lb.Reset("alice")

	if got := lb.Top(2); got != 60 {
		t.Fatalf("expected 60, got %d", got)
	}

	lb.AddScore("alice", 15)
	if got := lb.Top(2); got != 75 {
		t.Fatalf("expected 75, got %d", got)
	}
}

func TestTopWithNonPositiveK(t *testing.T) {
	lb := NewLeaderboard()
	lb.AddScore("alice", 100)

	if got := lb.Top(0); got != 0 {
		t.Fatalf("expected 0 for K=0, got %d", got)
	}
	if got := lb.Top(-1); got != 0 {
		t.Fatalf("expected 0 for K<0, got %d", got)
	}
}
