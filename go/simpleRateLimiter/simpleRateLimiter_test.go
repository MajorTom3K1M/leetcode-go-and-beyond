package main

import (
	"testing"
	"time"
)

func TestRateLimiterBasic(t *testing.T) {
	rl := NewRateLimiter(2)
	base := time.Unix(0, 0)

	user := "alice"

	// 1st request: allowed
	if !rl.Allow(user, base) {
		t.Fatal("expected first request to be allowed")
	}

	// 2nd request: allowed
	if !rl.Allow(user, base.Add(10*time.Second)) {
		t.Fatal("expected second request to be allowed")
	}

	// 3rd request within same minute: blocked
	if rl.Allow(user, base.Add(20*time.Second)) {
		t.Fatal("expected third request to be blocked")
	}

	// after 70 seconds, old ones expired
	if !rl.Allow(user, base.Add(70*time.Second)) {
		t.Fatal("expected request after 70s to be allowed")
	}
}
