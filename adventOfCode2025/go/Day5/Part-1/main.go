package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Range struct {
	start int
	end   int
}

func countValidIds() int {
	file, err := os.Open("./adventOfCode2025/go/Day5/Part-1/input.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	validIds := 0
	validIdsRanges := []Range{}

	isRange := true
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			isRange = false
			continue
		}

		if isRange {
			ranges := strings.Split(line, "-")
			if len(ranges) != 2 {
				continue
			}

			start, err := strconv.Atoi(ranges[0])
			end, err := strconv.Atoi(ranges[1])

			if err != nil {
				continue
			}

			validIdsRanges = append(validIdsRanges, Range{start: start, end: end})
		} else {
			for _, validIdRange := range validIdsRanges {
				id, err := strconv.Atoi(line)
				if err != nil {
					continue
				}
				if id >= validIdRange.start && id <= validIdRange.end {
					validIds++
					break
				}
			}
		}

	}

	return validIds
}

func main() {
	result := countValidIds()
	fmt.Printf("Total Valid IDs: %d", result)
}
