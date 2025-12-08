package main

import (
	"bufio"
	"log"
	"os"
)

func countTachyonBeams() int {
	file, err := os.Open("./adventOfCode2025/go/Day7/Part-1/input.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	manifoldSlices := [][]string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		slice := []string{}
		for _, char := range line {
			slice = append(slice, string(char))
		}
		manifoldSlices = append(manifoldSlices, slice)
	}

	countSplit := 0
	for row := 0; row < len(manifoldSlices)-1; row++ {
		for col := 0; col < len(manifoldSlices[row]); col++ {
			current := manifoldSlices[row][col]
			down := manifoldSlices[row+1][col]

			if current == "S" || current == "|" {
				if down == "." {
					manifoldSlices[row+1][col] = "|"
				} else if down == "^" {
					countSplit++
					if col-1 >= 0 {
						manifoldSlices[row+1][col-1] = "|"

					}

					if col+1 < len(manifoldSlices[row]) {
						manifoldSlices[row+1][col+1] = "|"
					}
				}
			}
		}
	}

	// Print manifoldSlices for visualization
	// for _, slice := range manifoldSlices {
	// 	line := ""
	// 	for _, char := range slice {
	// 		line += char
	// 	}
	// 	log.Println(line)
	// }

	return countSplit
}

func main() {
	result := countTachyonBeams()
	log.Printf("Total Tachyon Beams: %d", result)
}
