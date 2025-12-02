package main

import (
	"fmt"
	"strconv"
)

func calPoints(operations []string) int {
	stack := []string{}
	for _, operation := range operations {
		if operation == "+" {
			if len(stack) > 1 {
				top1, _ := strconv.Atoi(stack[len(stack)-1])
				stack = stack[:len(stack)-1]

				top2, _ := strconv.Atoi(stack[len(stack)-1])

				sum := top1 + top2
				stack = append(stack, strconv.Itoa(top1))

				stack = append(stack, strconv.Itoa(sum))
			}
		} else if operation == "D" {
			top, _ := strconv.Atoi(stack[len(stack)-1])
			stack = append(stack, strconv.Itoa(top*2))
		} else if operation == "C" {
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, operation)
		}
	}

	total := 0
	for _, value := range stack {
		valuInt, _ := strconv.Atoi(value)
		total += valuInt
	}

	return total
}

func main() {
	operations := []string{"5", "2", "C", "D", "+"}
	fmt.Println(calPoints(operations))
}
