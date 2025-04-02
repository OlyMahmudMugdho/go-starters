package utils

import (
	"encoding/json"
	"net/http"
)

func Json(w http.ResponseWriter, message interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(message)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
