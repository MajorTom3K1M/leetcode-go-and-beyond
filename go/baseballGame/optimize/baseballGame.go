package main

import "strconv"

func calPoints(operations []string) int {
	sum := 0
	records := make([]int, 0, len(operations))

	for _, operation := range operations {
		if operation == "C" {
			records = records[:len(records)-1]
		} else if operation == "D" {
			records = append(records, records[len(records)-1]*2)
		} else if operation == "+" {
			records = append(records, records[len(records)-1]+records[len(records)-2])
		} else {
			num, _ := strconv.Atoi(operation)
			records = append(records, num)
		}
	}

	for _, num := range records {
		sum += num
	}

	return sum
}
