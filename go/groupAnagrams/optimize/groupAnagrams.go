package main

import "fmt"

func groupAnagrams(strs []string) [][]string {
	var anagramMap = make(map[[27]int][]string)

	for _, str := range strs {
		anagramKey := [27]int{}
		for _, r := range str {
			alphabetOrder := int(r - 'a' + 1)
			anagramKey[alphabetOrder]++
		}
		anagramMap[anagramKey] = append(anagramMap[anagramKey], str)
	}

	var groups [][]string
	for _, val := range anagramMap {
		groups = append(groups, val)
	}

	return groups
}

func main() {
	// strs := []string{"eat", "tea", "tan", "ate", "nat", "bat"}
	// strs := []string{"ddddddddddg", "dgggggggggg"}
	strs := []string{"tin", "ram", "zip", "cry", "pus", "jon", "zip", "pyx"}
	fmt.Println(groupAnagrams(strs))
}
