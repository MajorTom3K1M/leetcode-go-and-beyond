package main

import (
	"bufio"
	"log"
	"math"
	"os"
)

func totalOutputJoltage() int {
	file, err := os.Open("./adventOfCode2025/go/Day3/Part-2/input.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	sum := 0

	for scanner.Scan() {
		line := scanner.Text()
		start := 0
		end := len(line)
		joltageTotal := 0
		for i := 11; i >= 0; i-- {
			maxJoltage := 0
			maxJoltageIndex := -1
			for j, char := range line[start : end-i] {
				singleJoltage := int(char - '0')
				if singleJoltage > maxJoltage {
					maxJoltage = singleJoltage
					maxJoltageIndex = start + j
				}
			}
			start = maxJoltageIndex + 1
			joltageTotal += maxJoltage * int(math.Pow(10, float64(i)))

		}

		sum += joltageTotal
	}

	return sum
}

func main() {
	result := totalOutputJoltage()
	println("Total output joltage is:", result)
}
