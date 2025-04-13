package main

import (
	"advanced-middleware/middlewares"
	"advanced-middleware/utils"
	"fmt"
	"log"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func main() {
	// the second middleware executes first
	log.Println("server is running on port :8080")
	http.HandleFunc("/", utils.Chain(Hello, utils.Method("GET"), middlewares.Logging(), middlewares.SecondLogger()))
	http.ListenAndServe(":8080", nil)
}
