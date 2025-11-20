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

func groupAnagrams(strs []string) [][]string {
	var groups [][]string

	used := make(map[int]bool)
	for i, str := range strs {
		if _, ok := used[i]; ok {
			continue
		}

		group := []string{str}
		used[i] = true

		for j, tStr := range strs {
			if i != j && !used[j] && isAnagram(str, tStr) {
				group = append(group, tStr)
				used[j] = true
			}
		}

		groups = append(groups, group)
	}

	return groups
}

func main() {
	strs := []string{"eat", "tea", "tan", "ate", "nat", "bat"}
	fmt.Println(groupAnagrams(strs))
}
