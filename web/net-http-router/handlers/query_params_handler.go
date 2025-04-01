package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func QueryParamsHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("hello %v", name),
	})
}
