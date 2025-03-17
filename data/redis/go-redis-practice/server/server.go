package server

import (
	"go-redis-practice/routes"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Server struct {
	Port   string
	Router *mux.Router
}

func NewServer() *Server {
	return &Server{
		Port:   os.Getenv("SERVER_PORT"),
		Router: routes.RegisterRoutes(),
	}
}

func (s *Server) Start() {
	log.Println("Server is running on port", s.Port)
	log.Fatal(http.ListenAndServe(":"+s.Port, s.Router))
}
