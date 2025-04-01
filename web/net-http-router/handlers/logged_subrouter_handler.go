package handlers

import (
	"encoding/json"
	"net/http"
)

func LoggedSubRouterHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"message": "logged subrouter handler",
	})
}
