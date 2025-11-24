package main

import "sync"

type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (sc *SafeCounter) Increment() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.count++
}

func (sc *SafeCounter) Value() int {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.count
}

func CountConcurrently(incrementPerWorker int, workers int) int {
	sc := SafeCounter{
		count: 0,
	}

	var wg sync.WaitGroup
	wg.Add(workers)

	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := 0; i < incrementPerWorker; i++ {
				sc.Increment()
			}
		}()
	}

	wg.Wait()
	return sc.Value()
}
func main() {
	val := CountConcurrently(1000, 20)
	println(val)
}
