package main

func isPerfectSquare(num int) bool {
	left, right := 1, num

	var mid int
	for left <= right {
		mid = (left + right) / 2

		square := mid * mid
		if square == num {
			return true
		} else if square < num {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return false
}

func main() {
	num := 16
	result := isPerfectSquare(num)
	println("Is", num, "a perfect square?", result)

	num = 14
	result = isPerfectSquare(num)
	println("Is", num, "a perfect square?", result)
}
