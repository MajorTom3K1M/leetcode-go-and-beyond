package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func countInvalidIDs() int {
	file, err := os.Open("./adventOfCode2025/go/Day2/Part-1/input.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	sum := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		pairs := strings.Split(line, ",")

		for _, pair := range pairs {
			ids := strings.Split(pair, "-")
			if len(ids) != 2 {
				continue
			}

			start, err := strconv.Atoi(ids[0])
			if err != nil {
				log.Fatalf("failed to convert start ID to integer: %s", err)
			}

			end, err := strconv.Atoi(ids[1])
			if err != nil {
				log.Fatalf("failed to convert end ID to integer: %s", err)
			}

			for id := start; id <= end; id++ {
				idLen := len(strconv.Itoa(id))
				if idLen%2 == 0 {
					mid := idLen / 2
					firstHalf := strconv.Itoa(id)[:mid]
					secondHalf := strconv.Itoa(id)[mid:]
					if firstHalf == secondHalf {
						sum += id
					}
				}
			}
		}
	}

	return sum
}

func main() {
	result := countInvalidIDs()
	fmt.Println("Result:", result)
}
