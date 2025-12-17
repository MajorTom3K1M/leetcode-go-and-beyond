package main

import (
	"fmt"
	"slices"
)

type Value struct {
	Value     string
	Timestamp int
}

type KVStore struct {
	cache map[string][]Value
}

func NewKVStore() *KVStore {
	return &KVStore{
		cache: make(map[string][]Value),
	}
}

func (kv *KVStore) Set(key string, value string, timestamp int) {
	newValue := &Value{
		Value:     value,
		Timestamp: timestamp,
	}

	kv.cache[key] = append(kv.cache[key], *newValue)
}

func (kv *KVStore) Get(key string, timestamp int) string {
	caches, exists := kv.cache[key]
	if !exists {
		return ""
	}

	copyCache := make([]Value, len(caches))
	copy(copyCache, caches)

	slices.SortFunc(copyCache, func(a, b Value) int {
		return b.Timestamp - a.Timestamp
	})

	for _, value := range copyCache {
		if value.Timestamp <= timestamp {
			return value.Value
		}
	}

	return ""
}

func main() {
	kv := NewKVStore()
	kv.Set("price", "10", 1)
	kv.Set("price", "12", 4)

	fmt.Println(kv.Get("price", 1))
	fmt.Println(kv.Get("price", 3))
	fmt.Println(kv.Get("price", 4))
	fmt.Println(kv.Get("price", 100))
	fmt.Println(kv.Get("unknown", 100))
}
