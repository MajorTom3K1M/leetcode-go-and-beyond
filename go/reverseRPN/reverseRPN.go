package main

import (
	"fmt"
	"strconv"
)

func evalRPN(tokens []string) int {
	stack := []int{}

	for _, token := range tokens {
		topIndex := len(stack)
		if val, err := strconv.Atoi(token); err == nil {
			stack = append(stack, val)
		} else {
			top1 := stack[topIndex-1]
			top2 := stack[topIndex-2]
			stack = stack[:topIndex-1]

			switch token {
			case "+":
				stack[len(stack)-1] = top1 + top2
			case "-":
				stack[len(stack)-1] = top2 - top1
			case "*":
				stack[len(stack)-1] = top1 * top2
			case "/":
				stack[len(stack)-1] = top2 / top1
			}
		}

	}

	if len(stack) == 1 {
		return stack[0]
	}

	return 0
}

func main() {
	result := evalRPN([]string{"2", "1", "+", "3", "*"})
	fmt.Println("Result:", result)

	result = evalRPN([]string{"4", "13", "5", "/", "+"})
	fmt.Println("Result:", result)

	result = evalRPN([]string{"10", "6", "9", "3", "+", "-11", "*", "/", "*", "17", "+", "5", "+"})
	fmt.Println("Result:", result)
}
