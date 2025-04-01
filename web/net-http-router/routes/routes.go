package routes

import (
	"log"
	"net-http-router/handlers"
	"net-http-router/middleware"
	"net/http"
)

func Serve() {
	router := http.NewServeMux()
	router.HandleFunc("GET /demo", handlers.SimpleHandler)
	router.HandleFunc("GET /demo/{id}", handlers.PathParamHandler)
	router.HandleFunc("POST /demo", handlers.PostBodyHandler)

	subrouter := http.NewServeMux()

	subrouter.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	loggedSubRouter := http.NewServeMux()
	loggedSubRouter.HandleFunc("GET /logger/logged", handlers.LoggedSubRouterHandler)
	router.Handle("/", middleware.Logger(loggedSubRouter))
	log.Println("server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
