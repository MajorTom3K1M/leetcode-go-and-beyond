package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func countValidIds() int {
	file, err := os.Open("./adventOfCode2025/go/Day5/Part-2/input.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	validIdsRanges := [][2]int{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break
		}

		ranges := strings.Split(line, "-")
		if len(ranges) != 2 {
			continue
		}

		start, err := strconv.Atoi(ranges[0])
		end, err := strconv.Atoi(ranges[1])

		if err != nil {
			continue
		}

		newRange := [2]int{start, end}
		insertPos := len(validIdsRanges)

		for i := len(validIdsRanges) - 1; i >= 0; i-- {
			if validIdsRanges[i][0] > start {
				insertPos = i
			} else {
				break
			}
		}

		validIdsRanges = append(validIdsRanges, [2]int{})
		copy(validIdsRanges[insertPos+1:], validIdsRanges[insertPos:])
		validIdsRanges[insertPos] = newRange
	}

	mergedRanges := [][2]int{validIdsRanges[0]}
	for i := 1; i < len(validIdsRanges); i++ {
		if len(mergedRanges) > 0 && mergedRanges[len(mergedRanges)-1][1] >= validIdsRanges[i][0] {
			if mergedRanges[len(mergedRanges)-1][1] < validIdsRanges[i][1] {
				mergedRanges[len(mergedRanges)-1][1] = validIdsRanges[i][1]
			}
		} else if len(mergedRanges) == 0 || mergedRanges[len(mergedRanges)-1][1] < validIdsRanges[i][0] {
			mergedRanges = append(mergedRanges, validIdsRanges[i])
		}
	}

	totalValidIds := 0
	for _, r := range mergedRanges {
		totalValidIds += (r[1] - r[0]) + 1
	}

	return totalValidIds
}

func main() {
	result := countValidIds()
	fmt.Printf("Total Valid IDs: %d", result)
}
