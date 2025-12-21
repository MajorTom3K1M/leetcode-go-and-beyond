package main

import "time"

type Request struct {
	Timestamp time.Time
}

type RateLimitter struct {
	request       map[string][]*Request // UserID -> []Req Timestamp
	MaxRequest    int
	WindowSeconds int
}

func NewRateLimiter(maxRequest int, windowSeconds int) *RateLimitter {
	return &RateLimitter{
		request:       make(map[string][]*Request),
		MaxRequest:    maxRequest,
		WindowSeconds: windowSeconds,
	}
}

func (rl *RateLimitter) IsAllowed(userID string) bool {
	now := time.Now()
	windowStart := now.Add(-time.Duration(rl.WindowSeconds) * time.Second)

	requests, exists := rl.request[userID]
	if !exists {
		rl.request[userID] = []*Request{{Timestamp: now}}
		return true
	}

	for len(requests) > 0 && requests[0].Timestamp.Before(windowStart) {
		requests[0] = nil
		requests = requests[1:]
	}
	rl.request[userID] = requests

	if len(requests) < rl.MaxRequest {
		requests = append(requests, &Request{Timestamp: now})
		rl.request[userID] = requests
		return true
	}

	return false
}

func main() {
	limiter := NewRateLimiter(5, 60)

	userID := "user123"
	for i := 0; i < 7; i++ {
		allowed := limiter.IsAllowed(userID)
		if allowed {
			println("Request", i+1, "allowed")
		} else {
			println("Request", i+1, "denied")
		}
		time.Sleep(10 * time.Second)
	}
}
