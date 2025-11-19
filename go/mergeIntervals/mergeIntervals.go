package main

import (
	"fmt"
	"sort"
)

func merge(intervals [][]int) [][]int {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	for i := 1; i < len(intervals); i++ {
		if intervals[i-1][1] >= intervals[i][0] {
			if intervals[i-1][1] <= intervals[i][1] {
				if intervals[i-1][0] <= intervals[i][0] {
					intervals[i][0] = intervals[i-1][0]
				}
				intervals = append(intervals[:i-1], intervals[i:]...)
				i--
			} else if intervals[i-1][1] > intervals[i][1] {
				if intervals[i-1][0] >= intervals[i][0] {
					intervals[i-1][0] = intervals[i][0]
				}
				intervals = append(intervals[:i], intervals[i+1:]...)
				i--
			}
		}
	}

	return intervals
}

func main() {
	intervals := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	// intervals := [][]int{{1, 4}, {4, 5}}
	// intervals := [][]int{{1, 4}, {0, 4}}
	// intervals := [][]int{{1, 4}, {0, 0}}
	fmt.Println(merge(intervals))
}
