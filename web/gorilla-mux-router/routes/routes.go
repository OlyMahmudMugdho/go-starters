package routes

import (
	"gorilla-mux-router/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.GetAllTodos).Methods("GET")
	return router
}
