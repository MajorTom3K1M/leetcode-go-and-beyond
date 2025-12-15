package leaderboardv2

import (
	"testing"
	"time"
)

func TestAddAndGetPlayer(t *testing.T) {
	lb := NewLeaderboard()

	// Add player successfully
	if !lb.AddPlayer("p1", "Alice") {
		t.Error("should add player successfully")
	}

	// Get player
	p := lb.GetPlayer("p1")
	if p == nil {
		t.Fatal("player should not be nil")
	}
	if p.Name != "Alice" || p.Score != 0 {
		t.Error("player data mismatch")
	}
	if p.GamesPlayed != 0 {
		t.Error("new player should have 0 games played")
	}
	if p.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}

	// Duplicate ID should fail
	if lb.AddPlayer("p1", "Bob") {
		t.Error("duplicate ID should fail")
	}

	// Empty name should fail
	if lb.AddPlayer("p2", "") {
		t.Error("empty name should fail")
	}

	// Non-existent player
	if lb.GetPlayer("p999") != nil {
		t.Error("non-existent player should return nil")
	}
}

func TestRemovePlayer(t *testing.T) {
	lb := NewLeaderboard()
	lb.AddPlayer("p1", "Alice")

	if !lb.RemovePlayer("p1") {
		t.Error("should remove existing player")
	}
	if lb.GetPlayer("p1") != nil {
		t.Error("removed player should not exist")
	}

	if lb.RemovePlayer("p999") {
		t.Error("removing non-existent should return false")
	}
}

func TestAddScore(t *testing.T) {
	lb := NewLeaderboard()
	lb.AddPlayer("p1", "Alice")

	// Add positive score
	newScore, ok := lb.AddScore("p1", 100)
	if !ok || newScore != 100 {
		t.Errorf("expected 100, got %d", newScore)
	}

	p := lb.GetPlayer("p1")
	if p.Score != 100 {
		t.Error("score should be 100")
	}
	if p.GamesPlayed != 1 {
		t.Error("games played should be 1")
	}
	if p.LastActive.IsZero() {
		t.Error("LastActive should be set")
	}

	// Add more score
	newScore, _ = lb.AddScore("p1", 50)
	if newScore != 150 {
		t.Errorf("expected 150, got %d", newScore)
	}
	if lb.GetPlayer("p1").GamesPlayed != 2 {
		t.Error("games played should be 2")
	}

	// Add negative score
	newScore, _ = lb.AddScore("p1", -30)
	if newScore != 120 {
		t.Errorf("expected 120, got %d", newScore)
	}

	// Non-existent player
	_, ok = lb.AddScore("p999", 100)
	if ok {
		t.Error("non-existent player should return false")
	}
}

func TestSetScore(t *testing.T) {
	lb := NewLeaderboard()
	lb.AddPlayer("p1", "Alice")
	lb.AddScore("p1", 100)

	// Set score (correction)
	if !lb.SetScore("p1", 500) {
		t.Error("should set score")
	}

	p := lb.GetPlayer("p1")
	if p.Score != 500 {
		t.Error("score should be 500")
	}
	if p.GamesPlayed != 1 {
		t.Error("games played should NOT increment on SetScore")
	}

	// Non-existent player
	if lb.SetScore("p999", 100) {
		t.Error("non-existent player should return false")
	}
}

func TestGetRank(t *testing.T) {
	lb := NewLeaderboard()
	lb.AddPlayer("p1", "Alice")
	lb.AddPlayer("p2", "Bob")
	lb.AddPlayer("p3", "Charlie")

	lb.SetScore("p1", 100)
	lb.SetScore("p2", 200)
	lb.SetScore("p3", 150)

	// p2 has highest score = rank 1
	if rank := lb.GetRank("p2"); rank != 1 {
		t.Errorf("p2 should be rank 1, got %d", rank)
	}
	// p3 is second = rank 2
	if rank := lb.GetRank("p3"); rank != 2 {
		t.Errorf("p3 should be rank 2, got %d", rank)
	}
	// p1 is third = rank 3
	if rank := lb.GetRank("p1"); rank != 3 {
		t.Errorf("p1 should be rank 3, got %d", rank)
	}

	// Non-existent player
	if rank := lb.GetRank("p999"); rank != 0 {
		t.Error("non-existent player should return rank 0")
	}
}

func TestGetRankTies(t *testing.T) {
	lb := NewLeaderboard()
	lb.AddPlayer("p1", "Alice")
	lb.AddPlayer("p2", "Bob")
	lb.AddPlayer("p3", "Charlie")
	lb.AddPlayer("p4", "Diana")

	lb.SetScore("p1", 100)
	lb.SetScore("p2", 200)
	lb.SetScore("p3", 200) // tie with p2
	lb.SetScore("p4", 50)

	// p2 and p3 should both be rank 1
	if rank := lb.GetRank("p2"); rank != 1 {
		t.Errorf("p2 should be rank 1, got %d", rank)
	}
	if rank := lb.GetRank("p3"); rank != 1 {
		t.Errorf("p3 should be rank 1 (tie), got %d", rank)
	}
	// p1 should be rank 3 (not 2, because two players are above)
	if rank := lb.GetRank("p1"); rank != 3 {
		t.Errorf("p1 should be rank 3, got %d", rank)
	}
	// p4 should be rank 4
	if rank := lb.GetRank("p4"); rank != 4 {
		t.Errorf("p4 should be rank 4, got %d", rank)
	}
}

func TestGetTopN(t *testing.T) {
	lb := NewLeaderboard()
	lb.AddPlayer("p1", "Alice")
	lb.AddPlayer("p2", "Bob")
	lb.AddPlayer("p3", "Charlie")

	lb.SetScore("p1", 100)
	lb.SetScore("p2", 200)
	lb.SetScore("p3", 150)

	top := lb.GetTopN(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 players, got %d", len(top))
	}
	if top[0].ID != "p2" {
		t.Error("first should be p2 (highest score)")
	}
	if top[1].ID != "p3" {
		t.Error("second should be p3")
	}

	// Get more than available
	all := lb.GetTopN(100)
	if len(all) != 3 {
		t.Error("should return all 3 players")
	}

	// Get 0
	none := lb.GetTopN(0)
	if len(none) != 0 {
		t.Error("should return empty slice")
	}
}

func TestGetTopNTiebreaker(t *testing.T) {
	lb := NewLeaderboard()
	lb.AddPlayer("p1", "Charlie")
	lb.AddPlayer("p2", "Alice")
	lb.AddPlayer("p3", "Bob")

	// All same score
	lb.SetScore("p1", 100)
	lb.SetScore("p2", 100)
	lb.SetScore("p3", 100)

	top := lb.GetTopN(3)
	// Should be sorted by name: Alice, Bob, Charlie
	if top[0].Name != "Alice" || top[1].Name != "Bob" || top[2].Name != "Charlie" {
		t.Error("same score should sort by name ascending")
	}
}

func TestGetPlayersInRankRange(t *testing.T) {
	lb := NewLeaderboard()
	lb.AddPlayer("p1", "Alice")
	lb.AddPlayer("p2", "Bob")
	lb.AddPlayer("p3", "Charlie")
	lb.AddPlayer("p4", "Diana")
	lb.AddPlayer("p5", "Eve")

	lb.SetScore("p1", 500)
	lb.SetScore("p2", 400)
	lb.SetScore("p3", 300)
	lb.SetScore("p4", 200)
	lb.SetScore("p5", 100)

	// Get ranks 2-4
	players := lb.GetPlayersInRankRange(2, 4)
	if len(players) != 3 {
		t.Fatalf("expected 3 players, got %d", len(players))
	}
	if players[0].ID != "p2" || players[1].ID != "p3" || players[2].ID != "p4" {
		t.Error("should return players ranked 2, 3, 4")
	}

	// Invalid range
	players = lb.GetPlayersInRankRange(10, 20)
	if len(players) != 0 {
		t.Error("out of range should return empty")
	}

	// Partial range
	players = lb.GetPlayersInRankRange(4, 10)
	if len(players) != 2 {
		t.Errorf("should return 2 players (ranks 4 and 5), got %d", len(players))
	}
}

func TestGetPlayersAboveScore(t *testing.T) {
	lb := NewLeaderboard()
	lb.AddPlayer("p1", "Alice")
	lb.AddPlayer("p2", "Bob")
	lb.AddPlayer("p3", "Charlie")

	lb.SetScore("p1", 100)
	lb.SetScore("p2", 200)
	lb.SetScore("p3", 150)

	// Above 120
	players := lb.GetPlayersAboveScore(120)
	if len(players) != 2 {
		t.Fatalf("expected 2 players, got %d", len(players))
	}
	// Should be sorted by score descending
	if players[0].ID != "p2" || players[1].ID != "p3" {
		t.Error("should be sorted by score descending")
	}

	// Above 200 (exclusive)
	players = lb.GetPlayersAboveScore(200)
	if len(players) != 0 {
		t.Error("no player above 200")
	}

	// Above 0
	players = lb.GetPlayersAboveScore(0)
	if len(players) != 3 {
		t.Error("all players above 0")
	}
}

func TestGetScoreHistory(t *testing.T) {
	lb := NewLeaderboard()
	lb.AddPlayer("p1", "Alice")

	lb.AddScore("p1", 100)
	time.Sleep(10 * time.Millisecond)
	lb.AddScore("p1", 50)
	time.Sleep(10 * time.Millisecond)
	lb.SetScore("p1", 200)

	history := lb.GetScoreHistory("p1")
	if len(history) != 3 {
		t.Fatalf("expected 3 history entries, got %d", len(history))
	}

	// Check first entry
	if history[0].OldScore != 0 || history[0].NewScore != 100 || history[0].Change != 100 {
		t.Error("first entry incorrect")
	}
	// Check second entry
	if history[1].OldScore != 100 || history[1].NewScore != 150 || history[1].Change != 50 {
		t.Error("second entry incorrect")
	}
	// Check third entry (SetScore)
	if history[2].OldScore != 150 || history[2].NewScore != 200 || history[2].Change != 50 {
		t.Error("third entry incorrect")
	}

	// Should be sorted by timestamp
	for i := 1; i < len(history); i++ {
		if history[i].Timestamp.Before(history[i-1].Timestamp) {
			t.Error("history should be sorted by timestamp ascending")
		}
	}

	// Non-existent player
	if len(lb.GetScoreHistory("p999")) != 0 {
		t.Error("non-existent player should return empty slice")
	}
}

func TestGetRecentlyActive(t *testing.T) {
	lb := NewLeaderboard()
	lb.AddPlayer("p1", "Alice")
	lb.AddPlayer("p2", "Bob")
	lb.AddPlayer("p3", "Charlie")

	lb.AddScore("p1", 100)
	time.Sleep(50 * time.Millisecond)
	lb.AddScore("p2", 100)
	time.Sleep(50 * time.Millisecond)
	lb.AddScore("p3", 100)

	// Get active in last 70ms (should be p3 and p2)
	// Note: Using 70ms instead of 60ms to account for test execution overhead
	active := lb.GetRecentlyActive(70 * time.Millisecond)
	if len(active) != 2 {
		t.Fatalf("expected 2 active players, got %d", len(active))
	}
	// Should be sorted by LastActive descending (most recent first)
	if active[0].ID != "p3" {
		t.Error("p3 should be first (most recent)")
	}
	if active[1].ID != "p2" {
		t.Error("p2 should be second")
	}
}

func TestGetAverageScore(t *testing.T) {
	lb := NewLeaderboard()

	// Empty leaderboard
	if avg := lb.GetAverageScore(); avg != 0 {
		t.Error("empty leaderboard should return 0")
	}

	lb.AddPlayer("p1", "Alice")
	lb.AddPlayer("p2", "Bob")
	lb.AddPlayer("p3", "Charlie")

	lb.SetScore("p1", 100)
	lb.SetScore("p2", 200)
	lb.SetScore("p3", 300)

	avg := lb.GetAverageScore()
	if avg != 200.0 {
		t.Errorf("expected 200.0, got %f", avg)
	}
}

func TestResetAllScores(t *testing.T) {
	lb := NewLeaderboard()
	lb.AddPlayer("p1", "Alice")
	lb.AddPlayer("p2", "Bob")

	lb.SetScore("p1", 100)
	lb.SetScore("p2", 200)

	count := lb.ResetAllScores()
	if count != 2 {
		t.Errorf("expected 2 reset, got %d", count)
	}

	if lb.GetPlayer("p1").Score != 0 || lb.GetPlayer("p2").Score != 0 {
		t.Error("all scores should be 0")
	}

	// Check history recorded
	h1 := lb.GetScoreHistory("p1")
	lastEntry := h1[len(h1)-1]
	if lastEntry.NewScore != 0 {
		t.Error("reset should be recorded in history")
	}
}

func TestScoreHistoryMultiplePlayers(t *testing.T) {
	lb := NewLeaderboard()
	lb.AddPlayer("p1", "Alice")
	lb.AddPlayer("p2", "Bob")

	lb.AddScore("p1", 100)
	lb.AddScore("p2", 200)
	lb.AddScore("p1", 50)

	// Each player's history should only contain their updates
	h1 := lb.GetScoreHistory("p1")
	if len(h1) != 2 {
		t.Errorf("p1 should have 2 history entries, got %d", len(h1))
	}

	h2 := lb.GetScoreHistory("p2")
	if len(h2) != 1 {
		t.Errorf("p2 should have 1 history entry, got %d", len(h2))
	}
}

func TestNegativeScores(t *testing.T) {
	lb := NewLeaderboard()
	lb.AddPlayer("p1", "Alice")

	lb.AddScore("p1", -50)
	if lb.GetPlayer("p1").Score != -50 {
		t.Error("score can be negative")
	}

	lb.SetScore("p1", -100)
	if lb.GetPlayer("p1").Score != -100 {
		t.Error("score can be set to negative")
	}
}
