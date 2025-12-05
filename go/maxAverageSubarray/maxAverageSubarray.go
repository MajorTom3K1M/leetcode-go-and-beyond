package main

func findMaxAverage(nums []int, k int) float64 {
	n := len(nums)
	currentSum := 0

	for i := 0; i < k; i++ {
		currentSum += nums[i]
	}

	maxAvg := float64(currentSum) / float64(k)

	for i := k; i < n; i++ {
		currentSum += nums[i] - nums[i-k]
		avg := float64(currentSum) / float64(k)

		if avg > maxAvg {
			maxAvg = avg
		}
	}

	return maxAvg
}

func main() {
	nums := []int{1, 12, -5, -6, 50, 3}
	k := 4
	result := findMaxAverage(nums, k)
	println(result)
}
