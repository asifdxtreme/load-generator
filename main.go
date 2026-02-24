package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

type QueryResponse struct {
	Pod            string  `json:"pod"`
	ProcessingTime string  `json:"processing_time"`
	ProcessingMs   float64 `json:"processing_ms"`
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// CPU-intensive work: compute SHA-256 hash 500K iterations
	data := []byte("load-generator-seed")
	for i := 0; i < 500000; i++ {
		hash := sha256.Sum256(data)
		data = hash[:]
	}

	elapsed := time.Since(start)

	hostname, _ := os.Hostname()
	resp := QueryResponse{
		Pod:            hostname,
		ProcessingTime: elapsed.String(),
		ProcessingMs:   float64(elapsed.Milliseconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func runServer() {
	http.HandleFunc("/query", queryHandler)
	http.HandleFunc("/health", healthHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func runLoadTest(url string, concurrency int) {
	fmt.Printf("Load test: %d concurrent requests to %s\n\n", concurrency, url)

	var (
		mu        sync.Mutex
		wg        sync.WaitGroup
		durations []time.Duration
		errors    int
	)

	client := &http.Client{Timeout: 30 * time.Second}

	start := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reqStart := time.Now()
			resp, err := client.Get(url)
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
		fmt.Printf("All %d requests failed!\n", concurrency)
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

func main() {
	if len(os.Args) > 1 && os.Args[0+1] == "loadtest" {
		fs := flag.NewFlagSet("loadtest", flag.ExitOnError)
		url := fs.String("url", "http://load-generator/query", "Target URL to load test")
		concurrency := fs.Int("c", 10, "Number of concurrent requests")
		fs.Parse(os.Args[2:])
		runLoadTest(*url, *concurrency)
		return
	}

	runServer()
}
