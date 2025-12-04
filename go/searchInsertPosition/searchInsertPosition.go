package main

func searchInsert(nums []int, target int) int {
	left, right := 0, len(nums)-1

	var mid int
	for left <= right {
		mid := (left + right) / 2

		if nums[mid] == target {
			return mid
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	if nums[mid] < target {
		return mid + 1
	} else {
		return mid
	}
}

func main() {
	nums := []int{1, 3, 5, 6}
	target := 5
	result := searchInsert(nums, target)
	println("Result:", result)

	target = 2
	result = searchInsert(nums, target)
	println("Result:", result)
}
