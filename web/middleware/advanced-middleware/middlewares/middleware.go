package middlewares

import (
	middlewaretypes "advanced-middleware/middleware_types"
	"log"
	"net/http"
	"time"
)

// Logging logs all requests with its path and the time it took to process
func Logging() middlewaretypes.Middleware {

	// Create a new Middleware
	return func(next http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			log.Println("first logger...")
			start := time.Now()
			defer func() { log.Println(r.URL.Path, time.Since(start)) }()

			// Call the next middleware/handler in chain
			next(w, r)
		}
	}
}

// second logger just prints "seconf logger..." in the console
func SecondLogger() middlewaretypes.Middleware {

	// Create a new Middleware
	return func(next http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			log.Println("second logger...")

			// Call the next middleware/handler in chain
			next(w, r)
		}
	}
}
