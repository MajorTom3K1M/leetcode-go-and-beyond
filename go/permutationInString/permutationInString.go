package main

func checkInclusion(s1 string, s2 string) bool {
	n, m := len(s1), len(s2)

	if n > m {
		return false
	}

	s1Count := [26]int{}
	s2Count := [26]int{}

	for i := 0; i < n; i++ {
		s1Count[s1[i]-'a']++
		s2Count[s2[i]-'a']++
	}

	if s1Count == s2Count {
		return true
	}

	for i := n; i < m; i++ {
		s2Count[s2[i]-'a']++
		s2Count[s2[i-n]-'a']--

		if s1Count == s2Count {
			return true
		}
	}

	return false
}

func main() {
	s1 := "ab"
	s2 := "eidbaooo"
	println(checkInclusion(s1, s2))
}
