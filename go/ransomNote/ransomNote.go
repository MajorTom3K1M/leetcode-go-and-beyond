package main

import "fmt"

func canConstruct(ransomNote string, magazine string) bool {
	magazineMap := make(map[rune]int)
	ransomNoteMap := make(map[rune]int)

	for _, r := range magazine {
		magazineMap[r]++
	}

	for _, r := range ransomNote {
		ransomNoteMap[r]++
	}

	for key, val := range ransomNoteMap {
		magVal, ok := magazineMap[key]
		if !ok || magVal < val {
			return false
		}
	}

	return true
}

func main() {
	ransomNote := "aa"
	magazine := "aab"
	fmt.Println(canConstruct(ransomNote, magazine))
}
