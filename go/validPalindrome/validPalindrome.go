package main

import "unicode"

func isPalindrome(s string) bool {
	L := 0
	R := len(s) - 1
	for L <= R {
		for L < R && !isAlphanumeric(rune(s[L])) {
			L++
		}

		for L < R && !isAlphanumeric(rune(s[R])) {
			R--
		}

		if unicode.ToLower(rune(s[L])) != unicode.ToLower(rune(s[R])) {
			return false
		}

		L++
		R--
	}

	return true
}

func isAlphanumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func main() {
	s1 := "A man, a plan, a canal: Panama"
	println(isPalindrome(s1))
	s2 := "race a car"
	println(isPalindrome(s2))
}
