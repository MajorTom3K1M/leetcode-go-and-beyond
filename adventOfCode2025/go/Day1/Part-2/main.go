package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func wrapAroundDial(position int) int {
	rangeStart := 0
	rangeEnd := 100
	position = position % rangeEnd
	if position < rangeStart {
		position += rangeEnd
	}
	return position
}

func hitsZeroCount(direction string, start int, steps int) int {
	if steps == 0 {
		return 0
	}

	var distToZero int
	if direction == "L" {
		distToZero = start % 100
	} else if direction == "R" {
		distToZero = (100 - start) % 100
	}

	if distToZero == 0 {
		distToZero = 100
	}

	if steps < distToZero {
		return 0
	}

	remaining := steps - distToZero
	hits := 1 + (remaining / 100)
	return hits
}

func dialCalculator() int {
	file, err := os.Open("./adventOfCode2025/go/Day1/Part-2/input.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	dial := 50
	counter := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		direction := line[0:1]
		steps, err := strconv.Atoi(line[1:])
		if err != nil {
			log.Fatalf("failed to convert steps to integer: %s", err)
		}

		counter += hitsZeroCount(direction, dial, steps)

		if direction == "L" {
			dial -= steps
		} else if direction == "R" {
			dial += steps
		}

		dial = wrapAroundDial(dial)
	}

	return counter
}

func main() {
	result := dialCalculator()
	fmt.Println("Result:", result)
}
