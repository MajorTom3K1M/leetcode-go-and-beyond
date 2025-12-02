package main

import "fmt"

type Temperatures struct {
	temp  int
	index int
}

func dailyTemperatures(temperatures []int) []int {
	result := make([]int, len(temperatures))
	stack := []Temperatures{}

	for i, temp := range temperatures {
		for len(stack) > 0 && temp > stack[len(stack)-1].temp {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			result[top.index] = i - top.index
		}

		stack = append(stack, Temperatures{
			temp:  temp,
			index: i,
		})
	}

	return result
}

func main() {
	temperatures := []int{73, 74, 75, 71, 69, 72, 76, 73}
	result := dailyTemperatures(temperatures)
	fmt.Println("Result:", result)
}
