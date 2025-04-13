package utils

import (
	middlewaretypes "advanced-middleware/middleware_types"
	"net/http"
)

// Method ensures that url can only be requested with a specific method, else returns a 400 Bad Request
func Method(m string) middlewaretypes.Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...middlewaretypes.Middleware) http.HandlerFunc {

	for i := len(middlewares) - 1; i >= 0; i-- {
		f = middlewares[i](f)
	}
	return f

	/*
		// this code-block will call the last middleware first
		for _, m := range middlewares {
			f = m(f)
		}
		return f
	*/
}
