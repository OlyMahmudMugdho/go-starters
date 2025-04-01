package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func PathParamHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	name := r.URL.Query().Get("name")
	json.NewEncoder(w).Encode(map[string]string{
		"id":   id,
		"name": fmt.Sprintf("hello %v", name),
	})
}
