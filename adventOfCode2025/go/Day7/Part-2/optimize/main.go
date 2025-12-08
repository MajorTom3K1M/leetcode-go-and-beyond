package main

import (
	"bufio"
	"log"
	"os"
)

func countTachyonBeams() int {
	file, err := os.Open("./adventOfCode2025/go/Day7/Part-2/input.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	manifoldSlices := [][]string{}
	manifoldWays := [][]int{}

	sCol := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		slice := []string{}
		for i, char := range line {
			if char == 'S' {
				sCol = i
			}
			slice = append(slice, string(char))
		}
		manifoldSlices = append(manifoldSlices, slice)

		splitRow := make([]int, len(line))
		manifoldWays = append(manifoldWays, splitRow)
	}
	manifoldWays[0][sCol] = 1

	// more space efficient way without modifying manifoldSlices
	for row := 0; row < len(manifoldSlices)-1; row++ {
		for col := 0; col < len(manifoldSlices[row]); col++ {
			k := manifoldWays[row][col]
			if k == 0 {
				continue
			}

			down := manifoldSlices[row+1][col]

			switch down {
			case ".", "S":
				manifoldWays[row+1][col] += k

			case "^":
				if col-1 >= 0 {
					manifoldWays[row+1][col-1] += k
				}

				if col+1 < len(manifoldSlices[row]) {
					manifoldWays[row+1][col+1] += k
				}
			default:
				manifoldWays[row+1][col] += k
			}
		}
	}

	lastRow := manifoldWays[len(manifoldWays)-1]
	sum := 0
	for _, val := range lastRow {
		sum += val
	}

	return sum
}

func main() {
	result := countTachyonBeams()
	log.Printf("Total Timeline: %d", result)
}
