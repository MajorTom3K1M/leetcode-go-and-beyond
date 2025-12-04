package main

import (
	"bufio"
	"log"
	"os"
)

func totalOutputJoltage() int {
	file, err := os.Open("./adventOfCode2025/go/Day3/Part-1/input.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	sum := 0

	for scanner.Scan() {
		line := scanner.Text()
		max1Idx, max2Idx := -1, -1
		max1Val, max2Val := -1, -1

		for i, char := range line {
			joltage := int(char - '0')

			if max1Idx == -1 || joltage > max1Val {
				max1Idx = i
				max1Val = joltage
			}
		}

		if max1Idx != len(line)-1 {
			for i, char := range line[max1Idx+1:] {
				joltage := int(char - '0')
				if max2Idx == -1 || joltage > max2Val {
					max2Idx = i + max1Idx + 1
					max2Val = joltage
				}
			}
		} else {
			for i, char := range line[:max1Idx] {
				joltage := int(char - '0')
				if max2Idx == -1 || joltage > max2Val {
					max2Idx = i
					max2Val = joltage
				}
			}

			max1Val, max2Val = max2Val, max1Val
			max1Idx, max2Idx = max2Idx, max1Idx
		}

		sum += max1Val*10 + max2Val
	}

	return sum
}

func main() {
	result := totalOutputJoltage()
	println("Total output joltage is:", result)
}
