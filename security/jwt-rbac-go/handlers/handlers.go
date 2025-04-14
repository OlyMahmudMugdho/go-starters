package handlers

import (
	"fmt"
	"net/http"

	"jwt-rbac-go/auth"
	"jwt-rbac-go/middlewares"
)

func Public(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Public access")
}

func User(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middlewares.ContextKey("user")).(*auth.Claims)
	fmt.Fprintf(w, "Hello %s! You are a %s.\n", claims.Username, claims.Role)
}

func Admin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome Admin!")
}
