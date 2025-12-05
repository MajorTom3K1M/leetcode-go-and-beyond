package main

import "fmt"

func minSubArrayLen(target int, nums []int) int {
	L := 0
	R := 0

	sum := 0
	minLen := len(nums) + 1

	for R < len(nums) {
		sum += nums[R]

		for sum >= target {
			if R-L+1 < minLen {
				minLen = R - L + 1
			}
			sum -= nums[L]
			L++
		}

		R++
	}

	if minLen == len(nums)+1 {
		return 0
	}

	return minLen
}

func main() {
	nums := []int{2, 3, 1, 2, 4, 3}
	target := 7
	result := minSubArrayLen(target, nums)
	fmt.Printf("Result: %d", result)
}
