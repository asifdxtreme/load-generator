package main

import (
	"flag"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

func main() {
	url := flag.String("url", "http://localhost:8080/query", "Target URL to load test")
	concurrency := flag.Int("c", 100, "Number of concurrent requests")
	flag.Parse()

	fmt.Printf("Load test: %d concurrent requests to %s\n\n", *concurrency, *url)

	var (
		mu       sync.Mutex
		wg       sync.WaitGroup
		durations []time.Duration
		errors    int
	)

	client := &http.Client{Timeout: 30 * time.Second}

	start := time.Now()

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reqStart := time.Now()
			resp, err := client.Get(*url)
			elapsed := time.Since(reqStart)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				errors++
				return
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				errors++
				return
			}

			durations = append(durations, elapsed)
		}()
	}

	wg.Wait()
	totalTime := time.Since(start)

	if len(durations) == 0 {
		fmt.Printf("All %d requests failed!\n", *concurrency)
		return
	}

	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })

	var sum time.Duration
	for _, d := range durations {
		sum += d
	}

	percentile := func(p float64) time.Duration {
		idx := int(float64(len(durations)) * p)
		if idx >= len(durations) {
			idx = len(durations) - 1
		}
		return durations[idx]
	}

	fmt.Println("========== Results ==========")
	fmt.Printf("Total time:       %s\n", totalTime)
	fmt.Printf("Requests:         %d successful, %d failed\n", len(durations), errors)
	fmt.Println("-----------------------------")
	fmt.Printf("Avg response:     %s\n", sum/time.Duration(len(durations)))
	fmt.Printf("Min response:     %s\n", durations[0])
	fmt.Printf("Max response:     %s\n", durations[len(durations)-1])
	fmt.Printf("P50 response:     %s\n", percentile(0.50))
	fmt.Printf("P95 response:     %s\n", percentile(0.95))
	fmt.Printf("P99 response:     %s\n", percentile(0.99))
	fmt.Println("=============================")
}
