package lru

import (
	"testing"
)

func TestNewLRUCache(t *testing.T) {
	c := NewLRUCache(3)
	if c == nil {
		t.Fatal("cache should not be nil")
	}
	if c.Len() != 0 {
		t.Error("new cache should be empty")
	}
}

func TestPutAndGet(t *testing.T) {
	c := NewLRUCache(3)

	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)

	val, ok := c.Get("a")
	if !ok || val != 1 {
		t.Errorf("expected 1, got %d", val)
	}

	val, ok = c.Get("b")
	if !ok || val != 2 {
		t.Errorf("expected 2, got %d", val)
	}

	// Non-existent key
	_, ok = c.Get("z")
	if ok {
		t.Error("should return false for non-existent key")
	}
}

func TestUpdate(t *testing.T) {
	c := NewLRUCache(3)

	c.Put("a", 1)
	c.Put("a", 100) // update

	val, _ := c.Get("a")
	if val != 100 {
		t.Error("value should be updated")
	}

	if c.Len() != 1 {
		t.Error("update should not increase length")
	}
}

func TestEviction(t *testing.T) {
	c := NewLRUCache(3)

	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)
	// Order: c, b, a (most to least recent)

	c.Put("d", 4) // should evict "a"

	if _, ok := c.Get("a"); ok {
		t.Error("'a' should be evicted")
	}

	if c.Len() != 3 {
		t.Error("length should still be 3")
	}

	// b, c, d should exist
	if _, ok := c.Get("b"); !ok {
		t.Error("'b' should exist")
	}
	if _, ok := c.Get("c"); !ok {
		t.Error("'c' should exist")
	}
	if _, ok := c.Get("d"); !ok {
		t.Error("'d' should exist")
	}
}

func TestGetUpdatesRecency(t *testing.T) {
	c := NewLRUCache(3)

	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)
	// Order: c, b, a

	c.Get("a") // access "a", making it most recent
	// Order: a, c, b

	c.Put("d", 4) // should evict "b" (now least recent)

	if _, ok := c.Get("b"); ok {
		t.Error("'b' should be evicted")
	}
	if _, ok := c.Get("a"); !ok {
		t.Error("'a' should exist (was accessed)")
	}
}

func TestPutUpdatesRecency(t *testing.T) {
	c := NewLRUCache(3)

	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)
	// Order: c, b, a

	c.Put("a", 100) // update "a", making it most recent
	// Order: a, c, b

	c.Put("d", 4) // should evict "b"

	if _, ok := c.Get("b"); ok {
		t.Error("'b' should be evicted")
	}
	if _, ok := c.Get("a"); !ok {
		t.Error("'a' should exist")
	}
}

func TestDelete(t *testing.T) {
	c := NewLRUCache(3)

	c.Put("a", 1)
	c.Put("b", 2)

	if !c.Delete("a") {
		t.Error("should return true for existing key")
	}
	if c.Delete("a") {
		t.Error("should return false for already deleted key")
	}
	if c.Delete("z") {
		t.Error("should return false for non-existent key")
	}

	if c.Len() != 1 {
		t.Error("length should be 1")
	}
}

func TestClear(t *testing.T) {
	c := NewLRUCache(3)

	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)

	c.Clear()

	if c.Len() != 0 {
		t.Error("cache should be empty")
	}
	if _, ok := c.Get("a"); ok {
		t.Error("'a' should not exist after clear")
	}
}

func TestContains(t *testing.T) {
	c := NewLRUCache(3)

	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)
	// Order: c, b, a

	// Contains should NOT update recency
	if !c.Contains("a") {
		t.Error("'a' should exist")
	}
	if c.Contains("z") {
		t.Error("'z' should not exist")
	}

	// "a" should still be least recent
	c.Put("d", 4) // should evict "a"
	if c.Contains("a") {
		t.Error("'a' should be evicted (Contains didn't update recency)")
	}
}

func TestPeek(t *testing.T) {
	c := NewLRUCache(3)

	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)
	// Order: c, b, a

	// Peek should NOT update recency
	val, ok := c.Peek("a")
	if !ok || val != 1 {
		t.Error("Peek should return value")
	}

	// "a" should still be least recent
	c.Put("d", 4) // should evict "a"
	if _, ok := c.Peek("a"); ok {
		t.Error("'a' should be evicted (Peek didn't update recency)")
	}

	// Peek non-existent
	_, ok = c.Peek("z")
	if ok {
		t.Error("should return false for non-existent key")
	}
}

func TestKeys(t *testing.T) {
	c := NewLRUCache(3)

	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)
	// Order: c, b, a (most to least recent)

	keys := c.Keys()
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "c" || keys[1] != "b" || keys[2] != "a" {
		t.Errorf("expected [c, b, a], got %v", keys)
	}

	// Access "a" to make it most recent
	c.Get("a")
	// Order: a, c, b

	keys = c.Keys()
	if keys[0] != "a" || keys[1] != "c" || keys[2] != "b" {
		t.Errorf("expected [a, c, b], got %v", keys)
	}
}

func TestCapacityOne(t *testing.T) {
	c := NewLRUCache(1)

	c.Put("a", 1)
	c.Put("b", 2) // evicts "a"

	if _, ok := c.Get("a"); ok {
		t.Error("'a' should be evicted")
	}
	if val, ok := c.Get("b"); !ok || val != 2 {
		t.Error("'b' should exist with value 2")
	}
}

func TestZeroCapacity(t *testing.T) {
	c := NewLRUCache(0) // should default to 1

	c.Put("a", 1)
	if c.Len() != 1 {
		t.Error("should have capacity of at least 1")
	}
}
