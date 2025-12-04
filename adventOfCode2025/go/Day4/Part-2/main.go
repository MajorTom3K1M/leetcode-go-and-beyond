package main

import (
	"bufio"
	"log"
	"os"
)

func totalAccessibleRolls() int {
	file, err := os.Open("./adventOfCode2025/go/Day4/Part-2/input.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var grid [][]byte
	for scanner.Scan() {
		line := scanner.Text()
		grid = append(grid, []byte(line))
	}

	adjacentDeltas := [][2]int{{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	isAccessible := true

	newGrid := make([][]byte, len(grid))
	for i := range grid {
		newGrid[i] = make([]byte, len(grid[i]))
		copy(newGrid[i], grid[i])
	}

	totalRemovedRolls := 0
	for isAccessible {
		totalAccessibleRolls := 0

		for i := 0; i < len(grid); i++ {
			for j := 0; j < len(grid[i]); j++ {
				currentCell := grid[i][j]

				if currentCell != '@' {
					continue
				}

				adjacentRolls := 0
				for _, delta := range adjacentDeltas {
					adjacentRow := i + delta[0]
					adjacentCol := j + delta[1]
					if adjacentRow >= 0 && adjacentRow < len(grid) && adjacentCol >= 0 && adjacentCol < len(grid[i]) {
						adjacentCell := grid[adjacentRow][adjacentCol]
						if adjacentCell == '@' {
							adjacentRolls++
						}
					}
				}

				if adjacentRolls < 4 {
					newGrid[i][j] = '.'
					totalAccessibleRolls++
				}
			}
		}

		totalRemovedRolls += totalAccessibleRolls

		if totalAccessibleRolls == 0 {
			isAccessible = false
		}

		copy(grid, newGrid)
	}

	return totalRemovedRolls
}

func main() {
	result := totalAccessibleRolls()
	println("Total accessible rolls is:", result)
}
