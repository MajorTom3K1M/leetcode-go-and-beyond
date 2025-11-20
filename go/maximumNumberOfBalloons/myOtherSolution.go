package main

func maxNumberOfBalloonsSolution(text string) int {
	balloonMap := map[rune]int{
		'b': 0,
		'a': 0,
		'l': 0,
		'o': 0,
		'n': 0,
	}

	for _, r := range text {
		if _, ok := balloonMap[r]; ok {
			balloonMap[r]++
		}
	}

	min := balloonMap['b']
	for k, v := range balloonMap {
		if k == 'l' || k == 'o' {
			if min > v/2 {
				min = v / 2
			}
		} else {
			if min > v {
				min = v
			}
		}
	}

	return min
}
