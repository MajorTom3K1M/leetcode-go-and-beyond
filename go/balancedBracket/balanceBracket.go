package main

import (
	"fmt"
)

func IsBalanced(s string) bool {
	pairs := map[rune]rune{
		'}': '{',
		']': '[',
		')': '(',
	}

	stack := []rune{}
	for _, r := range s {
		switch r {
		case '{', '[', '(':
			stack = append(stack, r)
		case '}', ']', ')':
			if len(stack) == 0 {
				return false
			}

			top := stack[len(stack)-1]
			if val, ok := pairs[r]; ok && top == val {
				stack = stack[:len(stack)-1]
			} else {
				return false
			}
		default:
			return false
		}
	}

	return len(stack) == 0
}

func main() {
	fmt.Println(IsBalanced("()[]{}"))
	fmt.Println(IsBalanced("(]"))
	fmt.Println(IsBalanced("([{}])"))
	fmt.Println(IsBalanced("("))
}
