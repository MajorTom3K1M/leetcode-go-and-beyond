package main

func containsDuplicate(nums []int) bool {
	duplicateMap := make(map[int]bool, len(nums))

	for _, num := range nums {
		if _, ok := duplicateMap[num]; ok {
			return true
		}

		duplicateMap[num] = true
	}

	return false
}
