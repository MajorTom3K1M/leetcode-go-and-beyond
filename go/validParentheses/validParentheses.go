package main

import "fmt"

func isValid(s string) bool {
	stack := []rune{}
	for _, r := range s {
		if r == '(' || r == '{' || r == '[' {
			stack = append(stack, r)
		} else {
			if len(stack) > 0 {
				top := stack[len(stack)-1]
				if (top == '(' && r == ')') || (top == '{' && r == '}') || (top == '[' && r == ']') {
					stack = stack[:len(stack)-1]
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	return len(stack) == 0
}

func main() {
	s := "["
	fmt.Println(isValid(s))
}
