package main

import (
	"context"
	"fmt"
	"sync"
)

func FindFirst(nums []int, target int, workers int) int {
	n := len(nums)
	if n == 0 || workers <= 0 {
		return -1
	}
	if workers > n {
		workers = n
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultCh := make(chan int, workers)
	var wg sync.WaitGroup

	chunkSize := (n + workers - 1) / workers
	for w := 0; w < workers; w++ {
		start := w * chunkSize
		if start >= n {
			break
		}
		end := start + chunkSize
		if end > n {
			end = n
		}

		wg.Add(1)

		go func(start, end int) {
			defer wg.Done()

			for i := start; i < end; i++ {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if nums[i] == target {
					select {
					case resultCh <- i:
						cancel()
					case <-ctx.Done():
					}
					return
				}
			}
		}(start, end)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	firstIndex := -1
	for idx := range resultCh {
		if firstIndex == -1 || idx < firstIndex {
			firstIndex = idx
		}
	}

	return firstIndex
}

func main() {
	nums := []int{5, 3, 7, 2, 9, 2, 10, 2}
	target := 2
	workers := 3

	idx := FindFirst(nums, target, workers)
	fmt.Printf("First index of %d is: %d\n", target, idx) // expect 3
}
