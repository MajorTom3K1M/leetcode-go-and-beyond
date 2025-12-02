package main

import "fmt"

func dailyTemperatures(temperatures []int) []int {
	stack := [][2]int{}
	answer := make([]int, len(temperatures))

	for i, temp := range temperatures {
		length := len(stack)
		if length > 0 {
			top := stack[length-1]
			if temp > top[0] {
				for j := length - 1; j >= 0; j-- {
					inTop := stack[len(stack)-1]
					if temp > inTop[0] {
						stack = stack[:len(stack)-1]
						answer[inTop[1]] = i - inTop[1]
					} else {
						break
					}
				}
			}
		}
		val := [2]int{temp, i}
		stack = append(stack, val)
	}

	return answer
}

func main() {
	temperatures := []int{73, 74, 75, 71, 69, 72, 76, 73}
	fmt.Println(dailyTemperatures(temperatures))
}
