package main

import "fmt"

func majorityElement(nums []int) int {
	elementMap := make(map[int]int)
	majorityValue := int(len(nums) / 2)
	for _, num := range nums {
		elementMap[num]++
		if val, ok := elementMap[num]; ok && val > majorityValue {
			return num
		}
	}

	return 0
}

func main() {
	nums := []int{2, 2, 1, 1, 1, 2, 2}
	fmt.Println(majorityElement(nums))
}
