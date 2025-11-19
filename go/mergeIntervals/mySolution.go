package main

import "sort"

func mergeSolution(intervals [][]int) [][]int {
	intervalLen := len(intervals)
	if intervalLen == 1 || intervalLen == 0 {
		return intervals
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	merged := [][]int{intervals[0]}
	for i := 1; i < intervalLen; i++ {
		if merged[len(merged)-1][1] >= intervals[i][0] {
			if merged[len(merged)-1][1] < intervals[i][1] {
				merged[len(merged)-1][1] = intervals[i][1]
			}
		} else {
			merged = append(merged, intervals[i])
		}
	}

	return merged
}

// func main() {
// 	intervals := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
// 	fmt.Println(test(intervals))
// }
