package cache

import (
	"slices"
	"time"
)

type CacheEntry struct {
	Key       string
	Value     interface{}
	CreatedAt time.Time
	ExpiresAt time.Time
	HitCount  int
}

type Cache struct {
	entries     map[string]*CacheEntry
	maxSize     int
	accessOrder []string
}

func NewCache(maxSize int) *Cache {
	return &Cache{
		entries: make(map[string]*CacheEntry),
		maxSize: maxSize,
	}
}

// Set adds or updates a key with a TTL (time-to-live).
// If TTL is 0, the entry never expires.
// If cache is full (at maxSize), remove the least recently accessed entry before adding.
// Returns true if a new entry was added, false if an existing entry was updated.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) bool {
	cache, exists := c.entries[key]
	if exists {
		cache.Value = value
		cache.CreatedAt = time.Now()
		if ttl > 0 {
			cache.ExpiresAt = cache.CreatedAt.Add(ttl)
		} else {
			cache.ExpiresAt = time.Time{}
		}
		return false
	}

	if c.maxSize > 0 && len(c.entries) >= c.maxSize {
		lruKey := c.accessOrder[0]
		delete(c.entries, lruKey)
		c.accessOrder = c.accessOrder[1:]
	}

	createdAt := time.Now()
	var expiredAt time.Time
	if ttl > 0 {
		expiredAt = createdAt.Add(ttl)
	} else {
		expiredAt = time.Time{}
	}
	newEntry := &CacheEntry{
		Key:       key,
		Value:     value,
		CreatedAt: createdAt,
		ExpiresAt: expiredAt,
		HitCount:  0,
	}
	c.entries[key] = newEntry
	c.accessOrder = append(c.accessOrder, key)
	return true
}

func removeKeyFromSlice(slice []string, key string) []string {
	for i, v := range slice {
		if v == key {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// Get retrieves a value by key.
// Returns (value, true) if found and not expired.
// Returns (nil, false) if not found or expired.
// If expired, the entry should be deleted.
// Increments HitCount on successful get.
func (c *Cache) Get(key string) (interface{}, bool) {
	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	if !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt) {
		delete(c.entries, key)
		c.accessOrder = removeKeyFromSlice(c.accessOrder, key)
		return nil, false
	}

	entry.HitCount++
	c.accessOrder = removeKeyFromSlice(c.accessOrder, key)
	c.accessOrder = append(c.accessOrder, key)
	return entry.Value, true
}

// Delete removes a key from the cache.
// Returns true if the key existed, false otherwise.
func (c *Cache) Delete(key string) bool {
	_, exists := c.entries[key]
	if !exists {
		return false
	}

	delete(c.entries, key)
	c.accessOrder = removeKeyFromSlice(c.accessOrder, key)
	return true
}

// Has checks if a key exists and is not expired.
// Does NOT increment HitCount.
// Deletes the entry if expired.
func (c *Cache) Has(key string) bool {
	entry, exists := c.entries[key]
	if !exists {
		return false
	}

	if !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt) {
		delete(c.entries, key)
		c.accessOrder = removeKeyFromSlice(c.accessOrder, key)
		return false
	}

	return true
}

// Size returns the number of entries in the cache (including expired ones).
func (c *Cache) Size() int {
	return len(c.entries)
}

// Keys returns all keys in the cache (including expired ones).
func (c *Cache) Keys() []string {
	keys := make([]string, 0, len(c.entries))
	for key := range c.entries {
		keys = append(keys, key)
	}
	return keys
}

// Clear removes all entries from the cache.
// Returns the number of entries removed.
func (c *Cache) Clear() int {
	removed := len(c.entries)
	c.entries = make(map[string]*CacheEntry)
	c.accessOrder = []string{}
	return removed
}

// CleanupExpired removes all expired entries.
// Returns the number of entries removed.
func (c *Cache) CleanupExpired() int {
	removed := 0
	now := time.Now()
	for key, entry := range c.entries {
		if !entry.ExpiresAt.IsZero() && now.After(entry.ExpiresAt) {
			delete(c.entries, key)
			c.accessOrder = removeKeyFromSlice(c.accessOrder, key)
			removed++
		}
	}
	return removed
}

// GetStats returns cache statistics.
// Returns: (totalEntries, totalHits, expiredCount)
// expiredCount = number of currently expired entries (without deleting them)
func (c *Cache) GetStats() (int, int, int) {
	totalEntries := len(c.entries)
	totalHits := 0
	expiredCount := 0
	now := time.Now()
	for _, entry := range c.entries {
		totalHits += entry.HitCount
		if !entry.ExpiresAt.IsZero() && now.After(entry.ExpiresAt) {
			expiredCount++
		}
	}
	return totalEntries, totalHits, expiredCount
}

// GetMostAccessed returns the top N entries by HitCount, sorted descending.
// Only includes non-expired entries.
// Returns fewer than N if there aren't enough non-expired entries.
func (c *Cache) GetMostAccessed(n int) []*CacheEntry {
	entries := make([]*CacheEntry, 0, n)
	for _, entry := range c.entries {
		if !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt) {
			continue
		}
		entries = append(entries, entry)
	}

	slices.SortFunc(entries, func(a, b *CacheEntry) int {
		return b.HitCount - a.HitCount
	})

	if len(entries) > n {
		return entries[:n]
	}

	return entries
}
