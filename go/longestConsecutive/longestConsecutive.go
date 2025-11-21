package main

import "fmt"

func longestConsecutive(nums []int) int {
	numSet := map[int]struct{}{}
	for _, num := range nums {
		numSet[num] = struct{}{}
	}

	longest := 0
	for num, _ := range numSet {
		consecutiveLen := 0
		if _, ok := numSet[num-1]; !ok {
			nextNum := num + 1
			consecutiveLen += 1
			_, nextNumExist := numSet[nextNum]
			for nextNumExist {
				nextNum += 1
				consecutiveLen += 1
				_, nextNumExist = numSet[nextNum]
			}
			longest = max(longest, consecutiveLen)
		}
	}

	return longest
}

func main() {
	nums := []int{100, 4, 200, 1, 3, 2}

	fmt.Println(longestConsecutive(nums))
}
