package main

import "fmt"

func maxNumberOfBalloons(text string) int {
	textMap := make(map[rune]int)
	balloonMap := map[rune]int{
		'b': 1,
		'a': 1,
		'l': 2,
		'o': 2,
		'n': 1,
	}
	for _, r := range text {
		textMap[r]++
	}

	count := 0
	isValid := true
	for isValid {
		for key, val := range balloonMap {
			textVal, ok := textMap[key]
			if ok {
				if textVal-val >= 0 {
					textMap[key] = textVal - val
					continue
				} else {
					isValid = false
				}
			} else {
				isValid = false
			}
		}
		if isValid {
			count++
		}
	}

	return count
}

func main() {
	text := "nlaebolko"
	fmt.Println(maxNumberOfBalloons(text))
}
