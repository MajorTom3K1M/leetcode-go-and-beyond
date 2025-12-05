package main

func characterReplacement(s string, k int) int {
	L := 0
	R := 0
	longest := 0

	alphabetCount := make([]int, 26)
	maxCount := 0

	for R < len(s) {
		alphabetCount[int(s[R]-'A')]++
		if alphabetCount[int(s[R]-'A')] > maxCount {
			maxCount = alphabetCount[int(s[R]-'A')]
		}

		for (R - L + 1 - maxCount) > k {
			alphabetCount[int(s[L]-'A')]--
			L++
		}

		w := R - L + 1
		if w > longest {
			longest = w
		}

		R++
	}

	return longest
}

func main() {
	result := characterReplacement("AABABBA", 1)
	println(result)
}
