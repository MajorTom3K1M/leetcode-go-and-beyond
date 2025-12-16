package lru

import "fmt"

type LRUCache struct {
	capacity    int
	accessOrder []string
	cache       map[string]int
}

func NewLRUCache(capacity int) *LRUCache {
	if capacity <= 0 {
		capacity = 1
	}

	return &LRUCache{
		capacity:    capacity,
		accessOrder: make([]string, 0),
		cache:       make(map[string]int),
	}
}

func (c *LRUCache) removeAccessOrder(key string) {
	for i, accessKey := range c.accessOrder {
		if accessKey == key {
			c.accessOrder = append(c.accessOrder[:i], c.accessOrder[i+1:]...)
			break
		}
	}
}

func (c *LRUCache) evictAccessOrder() string {
	evicted := c.accessOrder[0]
	c.accessOrder = c.accessOrder[1:]
	return evicted
}

func (c *LRUCache) pushAccessOrder(key string) {
	c.accessOrder = append(c.accessOrder, key)
}

// Get retrieves a value by key.
// Returns (value, true) if found, (0, false) if not found.
// Marks the key as recently used.
func (c *LRUCache) Get(key string) (int, bool) {
	fmt.Printf("MAP : %v \n", c.cache)
	if value, exists := c.cache[key]; exists {
		c.removeAccessOrder(key)
		c.pushAccessOrder(key)
		return value, true
	}
	return 0, false
}

// Put adds or updates a key-value pair.
// If cache is at capacity, evict the least recently used item first.
// Marks the key as recently used.
func (c *LRUCache) Put(key string, value int) {
	if _, exists := c.cache[key]; exists {
		c.cache[key] = value

		c.removeAccessOrder(key)

		c.pushAccessOrder(key)

		return
	}

	if len(c.cache) >= c.capacity {
		evicted := c.evictAccessOrder()
		delete(c.cache, evicted)
	}

	c.cache[key] = value
	c.pushAccessOrder(key)
}

// Delete removes a key from the cache.
// Returns true if key existed, false otherwise.
func (c *LRUCache) Delete(key string) bool {
	if _, exists := c.cache[key]; exists {
		delete(c.cache, key)
		c.removeAccessOrder(key)
		return true
	}
	return false
}

// Len returns the number of items in the cache.
func (c *LRUCache) Len() int {
	return len(c.cache)
}

// Clear removes all items from the cache.
func (c *LRUCache) Clear() {
	c.cache = make(map[string]int)
	c.accessOrder = make([]string, 0)
}

// Contains checks if a key exists WITHOUT marking it as recently used.
func (c *LRUCache) Contains(key string) bool {
	if _, exists := c.cache[key]; exists {
		return true
	}
	return false
}

// Peek gets a value WITHOUT marking it as recently used.
// Returns (value, true) if found, (0, false) if not found.
func (c *LRUCache) Peek(key string) (int, bool) {
	if value, exists := c.cache[key]; exists {
		return value, true
	}
	return 0, false
}

// Keys returns all keys in order from most recent to least recent.
func (c *LRUCache) Keys() []string {
	result := make([]string, 0)
	for i := len(c.accessOrder) - 1; i >= 0; i-- {
		result = append(result, c.accessOrder[i])
	}
	return result
}
