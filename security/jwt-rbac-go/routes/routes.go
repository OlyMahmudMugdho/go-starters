package routes

import (
	"net/http"

	"jwt-rbac-go/auth"
	"jwt-rbac-go/handlers"
	"jwt-rbac-go/middlewares"
)

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Auth
	mux.HandleFunc("/login", auth.LoginHandler)
	mux.HandleFunc("/register", auth.RegisterHandler)

	// Public
	mux.HandleFunc("/public", handlers.Public)

	// Protected
	mux.Handle("/user", middlewares.JWTMiddleware(http.HandlerFunc(handlers.User)))
	mux.Handle("/admin", middlewares.JWTMiddleware(middlewares.RequireRole("admin", handlers.Admin)))

	return mux
}
