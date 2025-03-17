package handlers

import (
	"encoding/json"
	"fmt"
	"go-redis-practice/types"
	"go-redis-practice/utils"
	"net/http"
	"strconv"
)

type HashHandler struct {
	HashOps *utils.HashOps
}

func NewHashHandler() *HashHandler {
	HashOps := utils.NewHashOps()

	return &HashHandler{
		HashOps: HashOps,
	}

}

func (h *HashHandler) HSet(w http.ResponseWriter, r *http.Request) {

	var body types.HashRequest

	json.NewDecoder(r.Body).Decode(&body)

	fmt.Println("body", body)

	w.Header().Add("Content-Type", "application/json")
	response, err := h.HashOps.HSet(body.Key, body.Value)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatInt(response, 10)))
}

func (h *HashHandler) HGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	field := r.URL.Query().Get("field")

	w.Header().Add("Content-Type", "application/json")

	response, err := h.HashOps.HGet(key, field)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func (h *HashHandler) HDel(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	field := r.URL.Query().Get("field")

	w.Header().Add("Content-Type", "application/json")

	response, err := h.HashOps.HDel(key, field)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatInt(response, 10)))
}
