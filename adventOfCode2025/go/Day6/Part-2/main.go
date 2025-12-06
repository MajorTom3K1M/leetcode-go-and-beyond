package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func readLinesAndOperators(path string) ([]string, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	isOperator := false

	scanner := bufio.NewScanner(f)
	lines := []string{}
	operators := []string{}
	for scanner.Scan() {
		line := scanner.Text()

		if line[0] == '+' || line[0] == '*' {
			isOperator = true
		}

		if !isOperator {
			lines = append(lines, line)
		} else {
			operatorsStr := strings.Fields(line)
			for _, operatorStr := range operatorsStr {
				if len(operatorStr) == 0 {
					continue
				}
				operators = append(operators, operatorStr)
			}
		}
	}
	return lines, operators, scanner.Err()
}

func solveCompactor() int {
	lines, operators, err := readLinesAndOperators("./adventOfCode2025/go/Day6/Part-2/input.txt")
	if err != nil {
		log.Fatalf("failed to read lines and operators: %s", err)
	}

	maxLen := len(lines[0])
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	delims := make([]bool, maxLen+1)
	for i := 0; i < maxLen; i++ {
		allUnderscore := true
		for _, line := range lines {
			if line[i] != ' ' {
				allUnderscore = false
				break
			}
		}
		delims[i] = allUnderscore
	}
	delims[len(delims)-1] = true

	digits := [][]string{}
	for _, line := range lines {
		startIdx := 0
		col := 0
		for i := 0; i <= maxLen; i++ {
			isSplit := delims[i]
			if isSplit && startIdx < i {
				split := line[startIdx:i]
				startIdx = i + 1
				if col >= len(digits) {
					digits = append(digits, []string{})
				}
				digits[col] = append(digits[col], split)
				col++
			}
		}
	}

	for i := 0; i < len(digits)/2; i++ {
		digits[i], digits[len(digits)-1-i] = digits[len(digits)-1-i], digits[i]
	}

	for i := 0; i < len(operators)/2; i++ {
		operators[i], operators[len(operators)-1-i] = operators[len(operators)-1-i], operators[i]
	}

	sum := 0
	for i, digitCol := range digits {
		if i >= len(operators) {
			log.Fatalf("not enough operators: have %d columns but only %d operators", len(digits), len(operators))
		}

		maxWidth := 0
		for _, str := range digitCol {
			if len(str) > maxWidth {
				maxWidth = len(str)
			}
		}

		numbers := []int{}
		for charPos := 0; charPos < maxWidth; charPos++ {
			numStr := ""
			for rowIdx := 0; rowIdx < len(digitCol); rowIdx++ {
				if charPos < len(digitCol[rowIdx]) && digitCol[rowIdx][charPos] != ' ' {
					numStr += string(digitCol[rowIdx][charPos])
				}
			}
			if len(strings.TrimSpace(numStr)) > 0 {
				converted, err := strconv.Atoi(strings.TrimSpace(numStr))
				if err != nil {
					log.Fatalf("failed to convert string '%s' to int: %s", numStr, err)
				}
				numbers = append(numbers, converted)
			}
		}

		result := 0
		if operators[i] == "+" {
			for _, num := range numbers {
				result += num
			}
		} else if operators[i] == "*" {
			result = 1
			for _, num := range numbers {
				result *= num
			}
		}
		sum += result
	}

	return sum
}

func main() {
	result := solveCompactor()
	println(result)
}
