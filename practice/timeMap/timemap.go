package main

import "sort"

type Value struct {
	Timestamp int
	val       string
}

type TimeMap struct {
	store map[string][]*Value // Key -> []Value with timestamp
}

func NewTimeMap() *TimeMap {
	return &TimeMap{
		store: make(map[string][]*Value),
	}
}

func (tm *TimeMap) Set(key string, value string, timestamp int) {
	val, ok := tm.store[key]
	if !ok {
		newValueArr := make([]*Value, 0)
		val = newValueArr
	}

	newValue := &Value{
		val:       value,
		Timestamp: timestamp,
	}
	val = append(val, newValue)
	tm.store[key] = val
}

func (tm *TimeMap) Get(key string, timestamp int) string {
	values, ok := tm.store[key]
	if !ok {
		return ""
	}

	i := sort.Search(len(values), func(i int) bool {
		return values[i].Timestamp > timestamp
	})

	if i == 0 {
		return ""
	}

	return values[i-1].val
}

func main() {
	tm := NewTimeMap()

	tm.Set("weather", "sunny", 1)
	tm.Set("weather", "cloudy", 4)
	tm.Set("weather", "rainy", 8)

	println(tm.Get("weather", 1))  // "sunny" (exact match)
	println(tm.Get("weather", 3))  // "sunny" (largest timestamp <= 3 is 1)
	println(tm.Get("weather", 5))  // "cloudy" (largest timestamp <= 5 is 4)
	println(tm.Get("weather", 8))  // "rainy" (exact match)
	println(tm.Get("weather", 10)) // "rainy" (largest timestamp <= 10 is 8)
	println(tm.Get("weather", 0))  // "" (no timestamp <= 0)
	println(tm.Get("temp", 5))     // "" (key doesn't exist)
}
