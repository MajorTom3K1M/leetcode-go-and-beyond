package main

import (
	"fmt"
	"strconv"
)

func calPoints(operations []string) int {
	stack := []int{}

	for _, ops := range operations {
		if val, err := strconv.Atoi(ops); err == nil {
			stack = append(stack, val)
		}

		if ops == "C" {
			stack = stack[:len(stack)-1]
		} else if ops == "D" {
			top := stack[len(stack)-1]
			stack = append(stack, top*2)
		} else if ops == "+" {
			prev := stack[len(stack)-1]
			prevPrev := stack[len(stack)-2]
			stack = append(stack, prev+prevPrev)
		}
	}

	result := 0
	for _, val := range stack {
		result += val
	}

	return result
}

func main() {
	operations := []string{"5", "2", "C", "D", "+"}
	result := calPoints(operations)
	fmt.Println("Result:", result)

	operations = []string{"5", "-2", "4", "C", "D", "9", "+", "+"}
	result = calPoints(operations)
	fmt.Println("Result:", result)

	operations = []string{"1", "C"}
	result = calPoints(operations)
	fmt.Println("Result:", result)
}
