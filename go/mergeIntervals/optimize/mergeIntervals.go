package main

import (
	"fmt"
	"sort"
)

func merge(intervals [][]int) [][]int {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	var merged [][]int = [][]int{intervals[0]}
	for i := 1; i < len(intervals); i++ {
		if len(merged) > 0 && merged[len(merged)-1][1] >= intervals[i][0] {
			if merged[len(merged)-1][1] < intervals[i][1] {
				merged[len(merged)-1][1] = intervals[i][1]
			}
		} else if len(merged) == 0 || merged[len(merged)-1][1] < intervals[i][0] {
			merged = append(merged, intervals[i])
		}
	}

	return merged
}

func main() {
	// intervals := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	// intervals := [][]int{{1, 4}, {0, 0}}
	intervals := [][]int{{1, 4}, {2, 3}}
	fmt.Println(merge(intervals))
}
