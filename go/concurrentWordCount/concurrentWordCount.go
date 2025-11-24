package main

import (
	"fmt"
	"strings"
	"sync"
)

func CountWords(docs []string, workers int) map[string]int {
	wordCount := make(map[string]int)
	jobs := make(chan string, workers)
	results := make(chan map[string]int)

	var wg sync.WaitGroup
	wg.Add(workers)

	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for doc := range jobs {
				localCount := make(map[string]int)
				words := strings.Fields(doc)
				for _, word := range words {
					_, ok := localCount[word]
					if !ok {
						localCount[word] = 0
					}
					localCount[word]++
				}

				results <- localCount
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		for _, doc := range docs {
			jobs <- doc
		}
		close(jobs)
	}()

	for localCount := range results {
		for word, count := range localCount {
			_, ok := wordCount[word]
			if !ok {
				wordCount[word] = 0
			}
			wordCount[word] += count
		}
	}

	return wordCount
}

func main() {
	docs := []string{
		"Go is fun and Go is fast",
		"I love concurrency in Go",
		"Fun fun fun with goroutines",
	}

	counts := CountWords(docs, 3)
	for w, c := range counts {
		fmt.Printf("%s: %d\n", w, c)
	}
}
