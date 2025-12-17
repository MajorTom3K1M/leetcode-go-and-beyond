package main

import (
	"fmt"
	"sort"
)

type Value struct {
	Value     string
	Timestamp int
}

type KVStore struct {
	cache map[string][]Value
}

func NewKVStore() *KVStore {
	return &KVStore{cache: make(map[string][]Value)}
}

func (kv *KVStore) Set(key, value string, timestamp int) {
	kv.cache[key] = append(kv.cache[key], Value{Value: value, Timestamp: timestamp})
}

func (kv *KVStore) Get(key string, timestamp int) string {
	arr, ok := kv.cache[key]
	if !ok || len(arr) == 0 {
		return ""
	}

	// Find first index i where arr[i].Timestamp > timestamp
	i := sort.Search(len(arr), func(i int) bool {
		return arr[i].Timestamp > timestamp
	})

	// i-1 is the last timestamp <= query
	if i == 0 {
		return ""
	}
	return arr[i-1].Value
}

func main() {
	kv := NewKVStore()
	kv.Set("price", "10", 1)
	kv.Set("price", "12", 4)

	fmt.Println(kv.Get("price", 1))   // 10
	fmt.Println(kv.Get("price", 3))   // 10
	fmt.Println(kv.Get("price", 4))   // 12
	fmt.Println(kv.Get("price", 100)) // 12
	fmt.Println(kv.Get("unknown", 100))
}
