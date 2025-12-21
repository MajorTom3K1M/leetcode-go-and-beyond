package main

import (
	"cmp"
	"fmt"
	"slices"
)

type PlayerScore struct {
	Name  string
	Score int
}

type Leaderboard struct {
	scores map[string]*PlayerScore
}

func NewLeaderboard() *Leaderboard {
	return &Leaderboard{
		scores: make(map[string]*PlayerScore),
	}
}

func (l *Leaderboard) AddScore(playerName string, score int) {
	if _, ok := l.scores[playerName]; ok {
		l.scores[playerName].Score += score
		return
	}

	l.scores[playerName] = &PlayerScore{
		Name:  playerName,
		Score: score,
	}
}

func (l *Leaderboard) GetTopPlayers(n int) []PlayerScore {
	list := make([]*PlayerScore, 0, len(l.scores))
	for _, ps := range l.scores {
		list = append(list, ps)
	}

	slices.SortFunc(list, func(a, b *PlayerScore) int {
		if b.Score != a.Score {
			return b.Score - a.Score
		}

		return cmp.Compare(b.Name, a.Name)
	})

	topK := []PlayerScore{}
	nLen := n
	if len(list) < n {
		nLen = len(list)
	}
	for i := 0; i < nLen; i++ {
		topK = append(topK, *list[i])
	}

	return topK
}

func (l *Leaderboard) GetPlayerRank(playerName string) int {
	list := make([]*PlayerScore, 0, len(l.scores))
	for _, ps := range l.scores {
		list = append(list, ps)
	}

	slices.SortFunc(list, func(a, b *PlayerScore) int {
		if b.Score != a.Score {
			return b.Score - a.Score
		}

		return cmp.Compare(b.Name, a.Name)
	})

	rank := 0
	prevScore := 0
	for i, p := range list {
		if i == 0 || p.Score != prevScore {
			rank = i + 1
			prevScore = p.Score
		}
		if p.Name == playerName {
			return rank
		}
	}

	return -1
}

func (l *Leaderboard) GetPlayerScore(playerName string) int {
	if ps, ok := l.scores[playerName]; ok {
		return ps.Score
	}
	return -1
}

func (l *Leaderboard) Reset() {
	for k := range l.scores {
		delete(l.scores, k)
	}
}

func main() {
	lb := NewLeaderboard()

	lb.AddScore("Alice", 100)
	lb.AddScore("Bob", 250)
	lb.AddScore("Charlie", 180)
	lb.AddScore("Alice", 90) // Alice now has 150

	fmt.Println(lb.GetPlayerScore("Alice")) // 150

	players := lb.GetTopPlayers(2) // [{Bob, 250}, {Charlie, 180}]
	for _, p := range players {
		fmt.Printf("Player: %s, Score: %d\n", p.Name, p.Score)
	}

	fmt.Println(lb.GetPlayerRank("Bob"))     // 1
	fmt.Println(lb.GetPlayerRank("Charlie")) // 2
	fmt.Println(lb.GetPlayerRank("Alice"))   // 3
	fmt.Println(lb.GetPlayerRank("Unknown")) // -1

	lb.Reset()
	fmt.Println(lb.GetTopPlayers(10)) // []
}
