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

func dialCalculator() int {
	file, err := os.Open("./adventOfCode2025/go/Day1/Part-1/input.txt")
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

		if direction == "L" {
			dial -= steps
			dial = wrapAroundDial(dial)
		} else if direction == "R" {
			dial += steps
			dial = wrapAroundDial(dial)
		}

		if dial == 0 {
			counter++
		}
	}

	return counter
}

func main() {
	result := dialCalculator()
	fmt.Println("Result:", result)
}
