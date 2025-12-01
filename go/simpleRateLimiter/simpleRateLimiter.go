package main

import (
	"time"
)

type RateLimiter struct {
	limit int
	// userID -> list of timestamps (when they made requests)
	requests map[string]map[string]int
}

func NewRateLimiter(limitPerMinute int) *RateLimiter {
	return &RateLimiter{
		limit:    limitPerMinute,
		requests: make(map[string]map[string]int),
	}
}

// Allow returns true if this request is allowed for this user at time `now`
// or false if the user exceeded the limit in the last minute.
func (r *RateLimiter) Allow(userID string, now time.Time) bool {
	layout := "2006-01-02 15 04"
	formattedTime := now.Format(layout)

	if _, ok := r.requests[userID]; !ok {
		r.requests[userID] = map[string]int{}
	}

	for key := range r.requests[userID] {
		if key != formattedTime {
			delete(r.requests[userID], key)
		}
	}

	r.requests[userID][formattedTime]++

	return r.requests[userID][formattedTime] <= r.limit
}

// func (r *RateLimiter) AllowSlidingWindow(userID string, now time.Time) bool {
// 	window := time.Minute
// 	threshold := now.Add(-window)

// 	// get user's history
// 	history := r.requests[userID]

// 	// step 1: drop timestamps older than `threshold`
// 	// find first index that is >= threshold
// 	idx := 0
// 	for idx < len(history) && history[idx].Before(threshold) {
// 		idx++
// 	}
// 	// keep only recent part
// 	history = history[idx:]

// 	// step 2: check if user already reached limit
// 	if len(history) >= r.limit {
// 		// update stored history (because we cleaned old ones)
// 		r.requests[userID] = history
// 		return false
// 	}

// 	// step 3: append current request
// 	history = append(history, now)
// 	r.requests[userID] = history

// 	return true
// }
