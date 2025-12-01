package main

import "fmt"

func reverseString(s []byte) {
	L := 0
	R := len(s) - 1
	for L <= R {
		s[L], s[R] = s[R], s[L]
		L++
		R--
	}
}

func main() {
	str := []byte{'h', 'e', 'l', 'l', 'o'}
	reverseString(str)
	fmt.Println(string(str))
	str = []byte{'H', 'a', 'n', 'n', 'a', 'h'}
	reverseString(str)
	fmt.Println(string(str))
}
