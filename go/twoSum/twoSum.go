package main

import "fmt"

func twoSum(nums []int, target int) []int {
	complementMap := make(map[int]int)
	var sumIndex []int
	for index, num := range nums {
		complement := target - num
		val, ok := complementMap[complement]
		if ok {
			sumIndex = []int{val, index}
			return sumIndex
		} else {
			complementMap[num] = index
		}
	}

	return sumIndex
}

func main() {
	nums := []int{3, 2, 4}
	target := 6
	fmt.Println(twoSum(nums, target))
}
