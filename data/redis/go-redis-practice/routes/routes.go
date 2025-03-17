package routes

import (
	"go-redis-practice/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()

	pingHandler := handlers.NewPingHandler()
	stringHandler := handlers.NewStringHandler()
	hashHandler := handlers.NewHashHandler()
	listHandler := handlers.NewListHandler()

	router.HandleFunc("/ping", pingHandler.DoPing).Methods("GET")

	router.HandleFunc("/string/set", stringHandler.Set).Methods("POST")
	router.HandleFunc("/string/get", stringHandler.Get).Methods("GET")
	router.HandleFunc("/string/delete", stringHandler.Delete).Methods("DELETE")

	router.HandleFunc("/hash/hset", hashHandler.HSet).Methods("POST")
	router.HandleFunc("/hash/hget", hashHandler.HGet).Methods("GET")
	router.HandleFunc("/hash/hdel", hashHandler.HDel).Methods("DELETE")

	router.HandleFunc("/list/lpush", listHandler.LPush).Methods("POST")
	router.HandleFunc("/list/lrange", listHandler.LRange).Methods("GET")
	router.HandleFunc("/list/lpop", listHandler.LPop).Methods("DELETE")
	router.HandleFunc("/list/rpush", listHandler.RPush).Methods("POST")
	router.HandleFunc("/list/rpop", listHandler.RPop).Methods("DELETE")
	router.HandleFunc("/list/llen", listHandler.LLen).Methods("GET")
	router.HandleFunc("/list/lindex", listHandler.LIndex).Methods("GET")

	return router
}
