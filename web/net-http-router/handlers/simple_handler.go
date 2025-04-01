package handlers

import (
	"encoding/json"
	"net/http"
)

func SimpleHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"message": "hello world",
	})
}
