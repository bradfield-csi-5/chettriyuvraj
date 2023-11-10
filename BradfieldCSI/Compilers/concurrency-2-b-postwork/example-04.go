package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var responses = []string{
	"200 OK",
	"402 Payment Required",
	"418 I'm a teapot",
}

func randomDelay(maxMillis int) time.Duration {
	return time.Duration(rand.Intn(maxMillis)) * time.Millisecond
}

func query(endpoint string) string {
	// Simulate querying the given endpoint
	delay := randomDelay(100)
	time.Sleep(delay)

	i := rand.Intn(len(responses))
	return responses[i]
}

// Query each of the mirrors in parallel and return the first
// response (this approach increases the amount of traffic but
// significantly improves "tail latency")
func parallelQuery(endpoints []string) {
	results := make(chan string, len(endpoints))
	var wg sync.WaitGroup

	for i := range endpoints {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results <- query(endpoints[i])
		}(i)
	}

	res := <-results
	fmt.Println(res)
	wg.Wait()
	close(results)
}

func main() {
	var endpoints = []string{
		"https://fakeurl.com/endpoint",
		"https://mirror1.com/endpoint",
		"https://mirror2.com/endpoint",
	}

	// Simulate long-running server process.
	// Hint: What will happen to the server's memory usage if it runs
	// continuously for a very long time?
	for {
		go parallelQuery(endpoints)
		delay := randomDelay(100)
		time.Sleep(delay)
	}
}
