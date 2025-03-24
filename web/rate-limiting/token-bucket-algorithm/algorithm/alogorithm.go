package algorithm

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type TokenBucket struct {
	mu         sync.Mutex
	tokens     int
	maxTokens  int
	refillRate time.Duration
}

// NewTokenBucket initializes a TokenBucket.
func NewTokenBucket(maxTokens int, refillRate time.Duration) *TokenBucket {
	tb := &TokenBucket{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
	}
	go tb.refillTokens()
	return tb
}

// refillTokens replenishes tokens at a fixed rate.
func (tb *TokenBucket) refillTokens() {
	ticker := time.NewTicker(tb.refillRate)
	for range ticker.C {
		tb.mu.Lock()
		if tb.tokens < tb.maxTokens {
			tb.tokens++
			log.Printf("Refilled token: now %d tokens available", tb.tokens)
		}
		tb.mu.Unlock()
	}
}

// AllowRequest checks if a request is allowed.
func (tb *TokenBucket) AllowRequest(clientIP string) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if tb.tokens > 0 {
		tb.tokens--
		log.Printf("Request from %s allowed - Tokens left: %d", clientIP, tb.tokens)
		return true
	}

	log.Printf("Request from %s denied - No tokens left", clientIP)
	return false
}

// RateLimitMiddleware applies rate limiting based on client IP.
func RateLimitMiddleware(limit int, refillRate time.Duration) func(http.Handler) http.Handler {
	buckets := make(map[string]*TokenBucket)
	var mu sync.Mutex

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract only the IP from the remote address.
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				log.Printf("error parsing remote address: %v", err)
				host = r.RemoteAddr // fallback to full address
			}

			mu.Lock()
			if _, exists := buckets[host]; !exists {
				buckets[host] = NewTokenBucket(limit, refillRate)
				log.Printf("Created new token bucket for %s", host)
			}
			tb := buckets[host]
			mu.Unlock()

			if tb.AllowRequest(host) {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "429 - Too Many Requests", http.StatusTooManyRequests)
			}
		})
	}
}
