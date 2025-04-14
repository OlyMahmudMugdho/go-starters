package main

import (
	"log"
	"net/http"

	"jwt-rbac-go/database"
	"jwt-rbac-go/routes"
)

func main() {
	database.Init()

	mux := routes.RegisterRoutes()

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
