package leaderboard

import "sort"

type Leaderboard struct {
	scores map[string]int
}

func NewLeaderboard() *Leaderboard {
	return &Leaderboard{
		scores: map[string]int{},
	}
}

// AddScore adds score to playerID total.
// If player does not exist, create them.
func (lb *Leaderboard) AddScore(playerID string, score int) {
	lb.scores[playerID] += score
}

// Top returns the sum of the top K scores.
// If there are fewer than K players, sum all of them.
// If K <= 0, return 0.
func (lb *Leaderboard) Top(K int) int {
	if K <= 0 {
		return 0
	}

	if len(lb.scores) == 0 {
		return 0
	}

	values := make([]int, 0, len(lb.scores))
	for _, s := range lb.scores {
		values = append(values, s)
	}

	sort.Ints(values)

	if K > len(values) {
		K = len(values)
	}

	sum := 0

	for i := len(values) - 1; i >= len(values)-K; i-- {
		sum += values[i]
	}

	return sum
}

// Reset resets the player's score to 0.
func (lb *Leaderboard) Reset(playerID string) {
	if _, ok := lb.scores[playerID]; ok {
		lb.scores[playerID] = 0
	}
}
