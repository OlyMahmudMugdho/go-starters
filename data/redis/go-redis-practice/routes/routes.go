package routes

import (
	"go-redis-practice/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()

	pingHandler := handlers.NewPingHandler()
	stringHandler := handlers.NewStringHandler()

	router.HandleFunc("/ping", pingHandler.DoPing).Methods("GET")

	router.HandleFunc("/string/set", stringHandler.Set).Methods("POST")
	router.HandleFunc("/string/get", stringHandler.Get).Methods("GET")
	router.HandleFunc("/string/delete", stringHandler.Delete).Methods("DELETE")

	return router
}
