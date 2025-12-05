package main

func longestOnes(nums []int, k int) int {
	L := 0
	R := 0
	zeroCount := 0
	maxLength := 0

	for R < len(nums) {
		if nums[R] == 0 {
			zeroCount++
		}

		for zeroCount > k {
			if nums[L] == 0 {
				zeroCount--
			}
			L++
		}

		if R-L+1 > maxLength {
			maxLength = R - L + 1
		}

		R++
	}

	return maxLength
}

func main() {
	nums := []int{1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0}
	k := 2
	result := longestOnes(nums, k)
	println(result)
}
