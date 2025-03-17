package routes

import (
	"go-redis-practice/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()

	pingHandler := handlers.NewPingHandler()

	router.HandleFunc("/ping", pingHandler.DoPing).Methods("GET")

	return router
}
