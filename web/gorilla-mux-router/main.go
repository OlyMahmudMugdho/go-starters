package main

import (
	"fmt"
	"gorilla-mux-router/routes"
	"log"
	"net/http"
)

const PORT int = 8080

func main() {
	router := routes.RegisterRoutes()
	log.Printf("server is listening on :%v", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", PORT), router))
}
