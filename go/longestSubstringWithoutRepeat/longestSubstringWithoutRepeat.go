package main

func lengthOfLongestSubstring(s string) int {
	L := 0
	R := 0
	maxWindow := 0

	duplicates := make(map[byte]struct{})

	for R < len(s) {
		_, found := duplicates[s[R]]
		for found {
			delete(duplicates, s[L])
			L++
			_, found = duplicates[s[R]]
		}

		duplicates[s[R]] = struct{}{}
		if currentWindow := R - L + 1; currentWindow > maxWindow {
			maxWindow = currentWindow
		}
		R++
	}

	return maxWindow
}

func main() {
	s := "abcabcbb"
	println(lengthOfLongestSubstring(s))
}
