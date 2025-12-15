package cache

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	c := NewCache(10)
	if c == nil {
		t.Fatal("cache should not be nil")
	}
	if c.Size() != 0 {
		t.Error("new cache should be empty")
	}
}

func TestSetAndGet(t *testing.T) {
	c := NewCache(10)

	// Set new key
	isNew := c.Set("key1", "value1", 0)
	if !isNew {
		t.Error("should return true for new key")
	}

	// Get existing key
	val, ok := c.Get("key1")
	if !ok {
		t.Error("should find key1")
	}
	if val != "value1" {
		t.Errorf("expected value1, got %v", val)
	}

	// Update existing key
	isNew = c.Set("key1", "value2", 0)
	if isNew {
		t.Error("should return false for update")
	}

	val, _ = c.Get("key1")
	if val != "value2" {
		t.Error("value should be updated")
	}

	// Get non-existent key
	_, ok = c.Get("key999")
	if ok {
		t.Error("should not find non-existent key")
	}
}

func TestExpiration(t *testing.T) {
	c := NewCache(10)

	// Set with short TTL
	c.Set("key1", "value1", 50*time.Millisecond)

	// Should exist immediately
	if _, ok := c.Get("key1"); !ok {
		t.Error("key should exist before expiration")
	}

	// Wait for expiration
	time.Sleep(60 * time.Millisecond)

	// Should be gone
	if _, ok := c.Get("key1"); ok {
		t.Error("key should be expired")
	}

	// Entry should be deleted after expired Get
	if c.Size() != 0 {
		t.Error("expired entry should be deleted")
	}
}

func TestHas(t *testing.T) {
	c := NewCache(10)
	c.Set("key1", "value1", 0)

	if !c.Has("key1") {
		t.Error("Has should return true for existing key")
	}
	if c.Has("key999") {
		t.Error("Has should return false for non-existent key")
	}

	// Has should not increment HitCount
	entry := c.entries["key1"]
	hitsBefore := entry.HitCount
	c.Has("key1")
	if entry.HitCount != hitsBefore {
		t.Error("Has should not increment HitCount")
	}

	// Get should increment HitCount
	c.Get("key1")
	if entry.HitCount != hitsBefore+1 {
		t.Error("Get should increment HitCount")
	}
}

func TestHasExpiration(t *testing.T) {
	c := NewCache(10)
	c.Set("key1", "value1", 50*time.Millisecond)

	time.Sleep(60 * time.Millisecond)

	if c.Has("key1") {
		t.Error("Has should return false for expired key")
	}
	if c.Size() != 0 {
		t.Error("Has should delete expired entry")
	}
}

func TestDelete(t *testing.T) {
	c := NewCache(10)
	c.Set("key1", "value1", 0)

	// Delete existing
	if !c.Delete("key1") {
		t.Error("should return true for existing key")
	}
	if c.Has("key1") {
		t.Error("key should be deleted")
	}

	// Delete non-existent
	if c.Delete("key999") {
		t.Error("should return false for non-existent key")
	}
}

func TestSizeAndKeys(t *testing.T) {
	c := NewCache(10)

	c.Set("key1", "value1", 0)
	c.Set("key2", "value2", 0)
	c.Set("key3", "value3", 0)

	if c.Size() != 3 {
		t.Errorf("expected size 3, got %d", c.Size())
	}

	keys := c.Keys()
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}

	// Check all keys present
	keyMap := make(map[string]bool)
	for _, k := range keys {
		keyMap[k] = true
	}
	if !keyMap["key1"] || !keyMap["key2"] || !keyMap["key3"] {
		t.Error("missing keys")
	}
}

func TestClear(t *testing.T) {
	c := NewCache(10)

	c.Set("key1", "value1", 0)
	c.Set("key2", "value2", 0)
	c.Set("key3", "value3", 0)

	removed := c.Clear()
	if removed != 3 {
		t.Errorf("expected 3 removed, got %d", removed)
	}
	if c.Size() != 0 {
		t.Error("cache should be empty after clear")
	}
}

func TestCleanupExpired(t *testing.T) {
	c := NewCache(10)

	c.Set("key1", "value1", 50*time.Millisecond)
	c.Set("key2", "value2", 50*time.Millisecond)
	c.Set("key3", "value3", 0) // never expires

	time.Sleep(60 * time.Millisecond)

	removed := c.CleanupExpired()
	if removed != 2 {
		t.Errorf("expected 2 removed, got %d", removed)
	}
	if c.Size() != 1 {
		t.Errorf("expected 1 remaining, got %d", c.Size())
	}
	if !c.Has("key3") {
		t.Error("key3 should still exist")
	}
}

func TestMaxSize(t *testing.T) {
	c := NewCache(3)

	c.Set("key1", "value1", 0)
	c.Set("key2", "value2", 0)
	c.Set("key3", "value3", 0)

	if c.Size() != 3 {
		t.Error("should have 3 entries")
	}

	// Access key1 and key3 to make key2 least recently accessed
	c.Get("key1")
	c.Get("key3")

	// Add new key - should evict key2 (least recently accessed)
	c.Set("key4", "value4", 0)

	if c.Size() != 3 {
		t.Errorf("size should still be 3, got %d", c.Size())
	}
	if c.Has("key2") {
		t.Error("key2 should be evicted (least recently accessed)")
	}
	if !c.Has("key1") || !c.Has("key3") || !c.Has("key4") {
		t.Error("key1, key3, key4 should exist")
	}
}

func TestMaxSizeUnlimited(t *testing.T) {
	c := NewCache(0) // unlimited

	for i := 0; i < 100; i++ {
		c.Set(string(rune('a'+i)), i, 0)
	}

	if c.Size() != 100 {
		t.Error("unlimited cache should hold all entries")
	}
}

func TestGetStats(t *testing.T) {
	c := NewCache(10)

	c.Set("key1", "value1", 0)
	c.Set("key2", "value2", 50*time.Millisecond)
	c.Set("key3", "value3", 0)

	// Access some keys
	c.Get("key1")
	c.Get("key1")
	c.Get("key3")

	time.Sleep(60 * time.Millisecond)

	total, hits, expired := c.GetStats()
	if total != 3 {
		t.Errorf("expected 3 total, got %d", total)
	}
	if hits != 3 {
		t.Errorf("expected 3 hits, got %d", hits)
	}
	if expired != 1 {
		t.Errorf("expected 1 expired, got %d", expired)
	}
}

func TestGetMostAccessed(t *testing.T) {
	c := NewCache(10)

	c.Set("key1", "value1", 0)
	c.Set("key2", "value2", 0)
	c.Set("key3", "value3", 0)
	c.Set("key4", "value4", 50*time.Millisecond) // will expire

	// Access with different frequencies
	c.Get("key1") // 1 hit
	c.Get("key2") // 3 hits
	c.Get("key2")
	c.Get("key2")
	c.Get("key3") // 2 hits
	c.Get("key3")
	c.Get("key4") // 5 hits but will expire
	c.Get("key4")
	c.Get("key4")
	c.Get("key4")
	c.Get("key4")

	time.Sleep(60 * time.Millisecond)

	// Get top 2 (should exclude expired key4)
	top := c.GetMostAccessed(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(top))
	}
	if top[0].Key != "key2" {
		t.Errorf("first should be key2, got %s", top[0].Key)
	}
	if top[1].Key != "key3" {
		t.Errorf("second should be key3, got %s", top[1].Key)
	}

	// Get more than available
	all := c.GetMostAccessed(100)
	if len(all) != 3 {
		t.Errorf("should return 3 non-expired entries, got %d", len(all))
	}
}

func TestSetUpdatesExpiration(t *testing.T) {
	c := NewCache(10)

	c.Set("key1", "value1", 50*time.Millisecond)
	time.Sleep(30 * time.Millisecond)

	// Update with new TTL
	c.Set("key1", "value2", 100*time.Millisecond)

	time.Sleep(40 * time.Millisecond)

	// Should still exist (new TTL from second Set)
	if !c.Has("key1") {
		t.Error("key should still exist after TTL update")
	}
}
