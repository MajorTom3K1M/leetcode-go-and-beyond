package leaderboardv2

import (
	"cmp"
	"math"
	"slices"
	"sync"
	"time"
)

type Player struct {
	ID          string
	Name        string
	Score       int
	GamesPlayed int
	LastActive  time.Time
	CreatedAt   time.Time
}

type ScoreUpdate struct {
	PlayerID  string
	OldScore  int
	NewScore  int
	Change    int
	Timestamp time.Time
}

type Leaderboard struct {
	players      map[string]*Player
	scoreHistory []ScoreUpdate
	mu           sync.Mutex
}

func NewLeaderboard() *Leaderboard {
	return &Leaderboard{
		players:      make(map[string]*Player),
		scoreHistory: make([]ScoreUpdate, 0),
	}
}

// AddPlayer adds a new player with initial score of 0.
// Returns false if player ID already exists or name is empty.
func (lb *Leaderboard) AddPlayer(id, name string) bool {
	if _, exists := lb.players[id]; exists {
		return false
	}

	if name == "" {
		return false
	}

	lb.players[id] = &Player{
		ID:          id,
		Name:        name,
		Score:       0,
		GamesPlayed: 0,
		CreatedAt:   time.Now(),
	}

	return true
}

// GetPlayer returns a player by ID, or nil if not found.
func (lb *Leaderboard) GetPlayer(id string) *Player {
	return lb.players[id]
}

// RemovePlayer removes a player by ID.
// Returns false if player doesn't exist.
func (lb *Leaderboard) RemovePlayer(id string) bool {
	if _, exists := lb.players[id]; exists {
		delete(lb.players, id)
		return true
	}
	return false
}

// AddScore adds points to a player's score (can be negative).
// Increments GamesPlayed by 1 and updates LastActive.
// Records the score change in history.
// Returns (newScore, true) if successful, (0, false) if player not found.
func (lb *Leaderboard) AddScore(playerID string, points int) (int, bool) {
	if player, exists := lb.players[playerID]; exists {
		oldScore := player.Score
		newScore := player.Score + points

		player.Score += points

		player.GamesPlayed++
		player.LastActive = time.Now()

		change := int(math.Abs(float64(oldScore - newScore)))

		lb.scoreHistory = append(lb.scoreHistory, ScoreUpdate{
			PlayerID:  playerID,
			OldScore:  oldScore,
			NewScore:  newScore,
			Change:    change,
			Timestamp: time.Now(),
		})

		return player.Score, true
	}
	return 0, false
}

// SetScore sets a player's score to an exact value.
// Does NOT increment GamesPlayed (used for corrections).
// Updates LastActive and records in history.
// Returns false if player not found.
func (lb *Leaderboard) SetScore(playerID string, score int) bool {
	if player, exists := lb.players[playerID]; exists {
		oldScore := player.Score
		newScore := score

		player.Score = newScore
		player.LastActive = time.Now()

		change := int(math.Abs(float64(oldScore - newScore)))

		lb.scoreHistory = append(lb.scoreHistory, ScoreUpdate{
			PlayerID:  playerID,
			OldScore:  oldScore,
			NewScore:  newScore,
			Change:    change,
			Timestamp: time.Now(),
		})

		return true
	}
	return false
}

// GetRank returns the rank of a player (1 = highest score).
// Players with the same score have the same rank.
// Returns 0 if player not found.
func (lb *Leaderboard) GetRank(playerID string) int {
	leaderboard := make([]*Player, 0, len(lb.players))
	for _, player := range lb.players {
		leaderboard = append(leaderboard, player)
	}

	slices.SortFunc(leaderboard, func(a, b *Player) int {
		if a.Score != b.Score {
			return b.Score - a.Score
		}
		return cmp.Compare(a.Name, b.Name)
	})

	rank := 0
	prevScore := 0
	for i, p := range leaderboard {
		if i == 0 || p.Score != prevScore {
			rank = i + 1
			prevScore = p.Score
		}
		if p.ID == playerID {
			return rank
		}
	}

	return 0
}

// GetTopN returns the top N players sorted by score descending.
// For same score, sort by name ascending (alphabetical).
// Returns fewer than N if there aren't enough players.
func (lb *Leaderboard) GetTopN(n int) []*Player {
	leaderboard := make([]*Player, 0, len(lb.players))
	for _, player := range lb.players {
		leaderboard = append(leaderboard, player)
	}

	slices.SortFunc(leaderboard, func(a, b *Player) int {
		if a.Score != b.Score {
			return b.Score - a.Score
		}
		return cmp.Compare(a.Name, b.Name)
	})

	end := n
	if len(leaderboard) < n {
		end = len(leaderboard)
	}

	return leaderboard[:end]
}

// GetPlayersInRankRange returns players whose rank is between startRank and endRank (inclusive).
// Sorted by rank ascending (highest score first).
// Example: GetPlayersInRankRange(1, 10) returns top 10 players.
func (lb *Leaderboard) GetPlayersInRankRange(startRank, endRank int) []*Player {
	leaderboard := make([]*Player, 0, len(lb.players))
	for _, player := range lb.players {
		leaderboard = append(leaderboard, player)
	}

	slices.SortFunc(leaderboard, func(a, b *Player) int {
		if a.Score != b.Score {
			return b.Score - a.Score
		}
		return cmp.Compare(a.Name, b.Name)
	})

	rank := 0
	prevScore := 0

	player := make([]*Player, 0)

	for i, p := range leaderboard {
		if i == 0 || prevScore != p.Score {
			rank = i + 1
			if rank >= startRank && rank <= endRank {
				player = append(player, p)
			}
		}
	}

	return player
}

// GetPlayersAboveScore returns all players with score > minScore.
// Sorted by score descending.
func (lb *Leaderboard) GetPlayersAboveScore(minScore int) []*Player {
	playerAboveScore := []*Player{}
	for _, player := range lb.players {
		if player.Score > minScore {
			playerAboveScore = append(playerAboveScore, player)
		}
	}

	slices.SortFunc(playerAboveScore, func(a, b *Player) int {
		if a.Score != b.Score {
			return b.Score - a.Score
		}
		return cmp.Compare(a.Name, b.Name)
	})

	return playerAboveScore
}

// GetScoreHistory returns all score updates for a player.
// Sorted by timestamp ascending (oldest first).
// Returns empty slice if player not found or no history.
func (lb *Leaderboard) GetScoreHistory(playerID string) []ScoreUpdate {
	history := []ScoreUpdate{}
	for _, update := range lb.scoreHistory {
		if update.PlayerID == playerID {
			history = append(history, update)
		}
	}

	slices.SortFunc(history, func(a, b ScoreUpdate) int {
		return a.Timestamp.Compare(b.Timestamp)
	})

	return history
}

// GetRecentlyActive returns players who were active within the given duration.
// Sorted by LastActive descending (most recent first).
func (lb *Leaderboard) GetRecentlyActive(within time.Duration) []*Player {
	cutoff := time.Now().Add(-within)
	activePlayers := []*Player{}
	for _, player := range lb.players {
		if player.LastActive.After(cutoff) || player.LastActive.Equal(cutoff) {
			activePlayers = append(activePlayers, player)
		}
	}

	slices.SortFunc(activePlayers, func(a, b *Player) int {
		return b.LastActive.Compare(a.LastActive)
	})

	return activePlayers
}

// GetAverageScore returns the average score of all players.
// Returns 0 if no players.
func (lb *Leaderboard) GetAverageScore() float64 {
	totalPlayer := len(lb.players)

	if totalPlayer == 0 {
		return 0
	}

	totalScore := 0
	for _, p := range lb.players {
		totalScore += p.Score
	}

	avg := float64(totalScore) / float64(totalPlayer)

	return avg
}

// ResetAllScores sets all players' scores to 0.
// Records each reset in history.
// Returns the number of players reset.
func (lb *Leaderboard) ResetAllScores() int {
	reset := len(lb.players)

	for _, p := range lb.players {
		p.Score = 0
	}

	lb.scoreHistory = []ScoreUpdate{}

	for _, p := range lb.players {
		lb.scoreHistory = append(lb.scoreHistory, ScoreUpdate{
			PlayerID:  p.ID,
			OldScore:  p.Score,
			NewScore:  0,
			Change:    p.Score,
			Timestamp: time.Now(),
		})
	}

	return reset
}
