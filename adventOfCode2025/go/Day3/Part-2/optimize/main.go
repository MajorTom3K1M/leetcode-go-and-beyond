package main

import (
	"bufio"
	"log"
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
		joltageTotal := maxNumberAfterRemoving12Digits(line)
		sum += joltageTotal
	}

	return sum
}

// maxNumberAfterRemoving12Digits finds the maximum number by removing exactly 12 digits
// Time complexity: O(n) where n is the length of the string
func maxNumberAfterRemoving12Digits(s string) int {
	n := len(s)
	toRemove := 12
	toKeep := n - toRemove

	// Use a stack to keep digits that form the maximum number
	stack := make([]byte, 0, toKeep)
	remainingRemoves := toRemove

	for i := 0; i < n; i++ {
		digit := s[i]

		// Remove smaller digits from stack if we can still remove digits
		// and if current digit is larger
		for len(stack) > 0 && remainingRemoves > 0 && stack[len(stack)-1] < digit {
			stack = stack[:len(stack)-1]
			remainingRemoves--
		}

		stack = append(stack, digit)
	}

	// If we still need to remove digits, remove from the end
	stack = stack[:toKeep]

	// Convert stack to integer
	result := 0
	for _, digit := range stack {
		result = result*10 + int(digit-'0')
	}

	return result
}

func main() {
	result := totalOutputJoltage()
	println("Total output joltage is:", result)
}
