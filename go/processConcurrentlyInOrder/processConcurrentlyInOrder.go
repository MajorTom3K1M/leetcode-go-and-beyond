package main

import (
	"fmt"
	"sync"
)

type job struct {
	index int
	value int
}

type result struct {
	index int
	value int
}

func square(n int) int {
	return n * n
}

func ProcessInts(nums []int, workers int) []int {
	jobs := make(chan job, workers)
	results := make(chan result, len(nums))

	var wg sync.WaitGroup
	wg.Add(workers)

	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()

			for job := range jobs {
				squaredValue := square(job.value)
				results <- result{index: job.index, value: squaredValue}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		for index, num := range nums {
			jobs <- job{index: index, value: num}
		}
		close(jobs)
	}()

	out := make([]int, len(nums))
	for result := range results {
		out[result.index] = result.value
	}

	return out
}

func main() {
	nums := []int{5, 1, 7, 3, 9}
	res := ProcessInts(nums, 3)
	fmt.Println("input :", nums)
	fmt.Println("output:", res)
}
