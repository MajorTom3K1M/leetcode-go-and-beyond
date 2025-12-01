package main

import "fmt"

func twoSum(numbers []int, target int) []int {
	L := 0
	R := len(numbers) - 1
	for L <= R {
		result := numbers[L] + numbers[R]
		if result > target {
			R--
		} else if result < target {
			L++
		} else if result == target {
			return []int{L, R}
		}
	}

	return []int{}
}

func main() {
	numbers := []int{2, 7, 11, 15}
	target := 9
	result := twoSum(numbers, target)
	fmt.Println(result)
}
