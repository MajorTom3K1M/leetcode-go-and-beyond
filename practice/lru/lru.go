package main

import "slices"

type LRUCache struct {
	store          map[string]int // KV store
	keyAccessOrder []string       // order of accessing key ascending
	capacity       int
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		store:          make(map[string]int),
		keyAccessOrder: make([]string, 0),
		capacity:       capacity,
	}
}

func (lru *LRUCache) RemoveAccessKey(key string) {
	lru.keyAccessOrder = slices.DeleteFunc(lru.keyAccessOrder, func(accessKey string) bool {
		return key == accessKey
	})
}

func (lru *LRUCache) CleanUpLeastUsed() {
	leastUsed := lru.keyAccessOrder[0]
	lru.RemoveAccessKey(leastUsed)
	delete(lru.store, leastUsed)
}

func (lru *LRUCache) Get(key string) (int, bool) {
	val, ok := lru.store[key]
	if !ok {
		return 0, false
	}

	lru.RemoveAccessKey(key)
	lru.keyAccessOrder = append(lru.keyAccessOrder, key)

	return val, true
}

func (lru *LRUCache) Put(key string, value int) {
	_, ok := lru.store[key]
	if ok {
		lru.RemoveAccessKey(key)
		lru.keyAccessOrder = append(lru.keyAccessOrder, key)
		lru.store[key] = value
		return
	}

	if len(lru.store) >= lru.capacity {
		lru.CleanUpLeastUsed()
		lru.keyAccessOrder = append(lru.keyAccessOrder, key)
		lru.store[key] = value
		return
	}

	lru.store[key] = value
	lru.keyAccessOrder = append(lru.keyAccessOrder, key)
}

func main() {
	cache := NewLRUCache(2)

	cache.Put("a", 1)
	cache.Put("b", 2)
	println(cache.Get("a")) // returns (1, true) - "a" is now most recent
	cache.Put("c", 3)       // evicts "b" (least recent)
	println(cache.Get("b")) // returns (0, false) - "b" was evicted
	println(cache.Get("a")) // returns (1, true)
	println(cache.Get("c")) // returns (3, true)
	cache.Put("d", 4)       // evicts "a" or "c"? (hint: "c" was used last)
	println(cache.Get("a"))
}
