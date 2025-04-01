package handlers

import (
	"encoding/json"
	"net-http-router/models"
	"net/http"
)

func PostBodyHandler(w http.ResponseWriter, r *http.Request) {
	var msg models.Message

	json.NewDecoder(r.Body).Decode(&msg)

	json.NewEncoder(w).Encode(map[string]models.Message{
		"data": msg,
	})
}
