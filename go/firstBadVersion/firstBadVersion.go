package main

/**
 * Forward declaration of isBadVersion API.
 * @param   version   your guess about first bad version
 * @return 	 	      true if current version is bad
 *			          false if current version is good
 * func isBadVersion(version int) bool;
 */
func isBadVersion(version int) bool {
	return version >= 4
}

func firstBadVersion(n int) int {
	left, right := 1, n

	var mid int
	for left < right {
		mid = left + (right-left)/2

		if isBadVersion(mid) {
			right = mid
		} else {
			left = mid + 1
		}
	}

	return left
}

func main() {
	n := 5
	result := firstBadVersion(n)
	println("First bad version is:", result)

	n = 1
	result = firstBadVersion(n)
	println("First bad version is:", result)
}
