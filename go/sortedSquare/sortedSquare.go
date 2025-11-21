package main

import "fmt"

func square(n int) int {
	return n * n
}

func sortedSquares(nums []int) []int {
	numsLen := len(nums)

	if numsLen == 0 {
		return nums
	}

	if numsLen == 1 {
		return []int{square(nums[0])}
	}

	result := make([]int, numsLen)
	L := 0
	R := numsLen - 1
	index := numsLen - 1
	for L <= R {
		if square(nums[R]) > square(nums[L]) {
			result[index] = square(nums[R])
			R--
		} else {
			result[index] = square(nums[L])
			L++
		}
		index--
	}

	return result
}

func main() {
	nums := []int{-4, -1, 0, 3, 10}

	fmt.Println(sortedSquares(nums))

	nums = []int{-7, -3, 2, 3, 11}

	fmt.Println(sortedSquares(nums))
}
