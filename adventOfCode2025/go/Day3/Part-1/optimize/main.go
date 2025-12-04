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
		maxJoltage := 0
		maxFirst := -1

		for i := 0; i < len(line); i++ {
			digit := int(line[i] - '0')

			if maxFirst >= 0 {
				joltage := maxFirst*10 + digit
				if joltage > maxJoltage {
					maxJoltage = joltage
				}
			}

			if digit > maxFirst {
				maxFirst = digit
			}
		}

		sum += maxJoltage
	}

	return sum
}

func main() {
	result := totalOutputJoltage()
	println("Total output joltage is:", result)
}

/*
Algorithm Explanation - O(n) Single Pass Solution
==================================================

Problem: Find the maximum 2-digit number that can be formed by selecting
two digits from a string while maintaining their original order.

Key Insight: As we scan left-to-right, we track the largest digit seen so far.
For each new digit, we try pairing it as the "second digit" with the best
"first digit" we've seen before it.

Example Walkthrough with "987654321111111":
-------------------------------------------
Step | Current | MaxFirst | Try Combination      | Best So Far
-----|---------|----------|---------------------|-------------
1    | 9       | none     | (skip - need 2)     | 0
2    | 8       | 9        | 9 + 8 = 98         | 98  ← Found it!
3    | 7       | 9        | 9 + 7 = 97         | 98  (no change)
4    | 6       | 9        | 9 + 6 = 96         | 98  (no change)
...  | ...     | 9        | All < 98           | 98

Result: 98 ✓

Example Walkthrough with "818181911112111":
-------------------------------------------
Step | Current | MaxFirst | Try Combination      | Best So Far
-----|---------|----------|---------------------|-------------
1    | 8       | none     | (skip - need 2)     | 0
2    | 1       | 8        | 8 + 1 = 81         | 81
3    | 8       | 8        | 8 + 8 = 88         | 88
4    | 1       | 8        | 8 + 1 = 81         | 88
5    | 8       | 8        | 8 + 8 = 88         | 88
6    | 1       | 8        | 8 + 1 = 81         | 88
7    | 9       | 8        | 8 + 9 = 89         | 89
8    | 1       | 9        | 9 + 1 = 91         | 91  ← 9 is now best first
9    | 1       | 9        | 9 + 1 = 91         | 91
10   | 1       | 9        | 9 + 1 = 91         | 91
11   | 2       | 9        | 9 + 2 = 92         | 92  ← Best answer!
12   | 1       | 9        | 9 + 1 = 91         | 92
...  | ...     | 9        | All < 92           | 92

Result: 92 ✓

Why This Works:
---------------
- We maintain order: first digit always comes before second digit
- We're greedy: always use the largest possible first digit
- Single pass: O(n) time complexity, O(1) space
- For any position j, the best 2-digit number ending there is:
  (max digit from positions 0 to j-1) * 10 + (digit at j)
*/
