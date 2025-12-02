package main

import "fmt"

func dailyTemperatures(temperatures []int) []int {
	stack := []int{}
	answer := make([]int, len(temperatures))

	for i, temp := range temperatures {
		for len(stack) > 0 && temp > temperatures[stack[len(stack)-1]] {
			top := stack[len(stack)-1]
			answer[top] = i - top
			stack = stack[:len(stack)-1] // pop
		}
		stack = append(stack, i)
	}

	return answer
}

func main() {
	temperatures := []int{73, 74, 75, 71, 69, 72, 76, 73}
	fmt.Println(dailyTemperatures(temperatures))
}
