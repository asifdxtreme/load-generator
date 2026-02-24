package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

func main() {
	http.HandleFunc("/query", queryHandler)
	http.HandleFunc("/health", healthHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
