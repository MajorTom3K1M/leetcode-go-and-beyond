package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func solveCompactor() int {
	file, err := os.Open("./adventOfCode2025/go/Day6/Part-1/input.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	isOperator := false
	scanner := bufio.NewScanner(file)

	digits := make([][]int, 0)
	operators := []string{}
	for scanner.Scan() {
		line := scanner.Text()

		if line[0] == '+' || line[0] == '*' {
			isOperator = true
		}

		if !isOperator {
			digitsStr := strings.Fields(line)
			for col, digitStr := range digitsStr {
				converted, err := strconv.Atoi(digitStr)
				if err != nil {
					log.Fatalf("failed to convert string to int: %s", err)
				}
				if col >= len(digits) {
					digits = append(digits, []int{})
				}
				digits[col] = append(digits[col], converted)
			}
		} else {
			operatorsStr := strings.Split(line, " ")
			for _, operatorStr := range operatorsStr {
				if len(operatorStr) == 0 {
					continue
				}
				operators = append(operators, operatorStr)
			}
		}
	}

	sum := 0
	for i, digit := range digits {
		result := 0
		for _, num := range digit {
			if operators[i] == "+" {
				result += num
			} else if operators[i] == "*" {
				if result == 0 {
					result = 1
				}
				result *= num
			}
		}
		sum += result
	}

	return sum
}

func main() {
	result := solveCompactor()
	fmt.Println("Result:", result)
}
