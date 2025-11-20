package main

func twoSumSolution(nums []int, target int) []int {
	compositionMap := map[int]int{}
	for index, value := range nums {
		targetIndex, ok := compositionMap[target-value]
		if ok {
			return []int{targetIndex, index}
		} else {
			compositionMap[value] = index
		}
	}

	return []int{}
}
