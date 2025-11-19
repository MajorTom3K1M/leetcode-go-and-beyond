package main

import "fmt"

func canConstruct(ransomNote string, magazine string) bool {
	magazineMap := make(map[rune]int)

	for _, r := range magazine {
		magazineMap[r]++
	}

	for _, r := range ransomNote {
		value, ok := magazineMap[r]
		if (ok && value <= 0) || !ok {
			return false
		}
		magazineMap[r]--
	}

	return true
}

func main() {
	ransomNote := "aa"
	magazine := "ab"
	fmt.Println(canConstruct(ransomNote, magazine))
}
