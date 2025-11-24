package main

import (
	"math"
	"sync"
)

func isPrime(num int) bool {
	if num < 2 {
		return false
	}

	if num == 2 {
		return true
	}

	if num%2 == 0 {
		return false
	}

	limit := int(math.Sqrt(float64(num)))
	for i := 3; i <= limit; i += 2 {
		if num%i == 0 {
			return false
		}
	}

	return true
}

func CountPrimesConcurrently(n int, workers int) int {
	jobs := make(chan int, workers)
	results := make(chan int)

	var wg sync.WaitGroup
	wg.Add(workers)

	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for num := range jobs {
				if isPrime(num) {
					results <- 1
				} else {
					results <- 0
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		for i := 2; i < n; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	count := 0
	for val := range results {
		count += val
	}

	return count
}

func main() {
	n := 100
	workers := 4
	primeCount := CountPrimesConcurrently(n, workers)
	println("Number of primes less than", n, "is:", primeCount)
}
