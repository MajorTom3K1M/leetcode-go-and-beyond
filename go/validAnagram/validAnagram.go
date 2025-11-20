package main

import "fmt"

func isAnagram(s string, t string) bool {
	tMap := make(map[rune]int)

	for _, r := range t {
		tMap[r]++
	}

	for _, r := range s {
		tMap[r]--
	}

	for _, val := range tMap {
		if val != 0 {
			return false
		}
	}

	return true
}

func main() {
	s := "rat"
	t := "car"
	fmt.Println(isAnagram(s, t))
}
