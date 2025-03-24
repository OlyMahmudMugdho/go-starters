package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"token-bucket-algorithm/algorithm"
)

// TokenBucket represents the rate limiter for a client.

func main() {
	mux := http.NewServeMux()

	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	// Apply rate limiting: 5 requests allowed, with 1 token refilled per second.
	rateLimitedHandler := algorithm.RateLimitMiddleware(5, 1*time.Second)(helloHandler)
	mux.Handle("/", rateLimitedHandler)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", mux)
}
