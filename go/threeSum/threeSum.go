package main

import "sort"

func threeSum(nums []int) [][]int {
	sort.Ints(nums)

	n := len(nums)
	result := [][]int{}

	for i, num := range nums {
		if num > 0 {
			break
		} else if i > 0 && nums[i-1] == nums[i] {
			continue
		}

		lo, hi := i+1, n-1

		for lo < hi {
			sum := nums[lo] + nums[hi] + num
			if sum == 0 {
				result = append(result, []int{nums[lo], nums[hi], num})
				lo++
				hi--
				for lo < hi && nums[lo] == nums[lo-1] {
					lo++
				}
				for lo < hi && nums[hi] == nums[hi+1] {
					hi--
				}
			} else if sum < 0 {
				lo++
			} else if sum > 0 {
				hi--
			}

		}
	}

	return result
}

func main() {
	nums := []int{-1, 0, 1, 2, -1, -4}
	result := threeSum(nums)
	for _, triplet := range result {
		println(triplet[0], triplet[1], triplet[2])
	}
}
