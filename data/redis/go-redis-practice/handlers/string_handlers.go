package handlers

import (
	"encoding/json"
	"go-redis-practice/types"
	"go-redis-practice/utils"
	"log"
	"net/http"
)

type StringHandler struct {
	StringOps *utils.StringOps
}

func NewStringHandler() *StringHandler {
	StringOps := utils.NewStringOps()

	return &StringHandler{
		StringOps: StringOps,
	}
}

func (s *StringHandler) Set(w http.ResponseWriter, r *http.Request) {
	var body types.StringOpsRequest

	json.NewDecoder(r.Body).Decode(&body)

	w.Header().Add("Content-Type", "application/json")

	response, err := s.StringOps.Set(body.Key, body.Value)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(&types.StringOpsSetResponse{
		Status: response,
	})
}

func (s *StringHandler) Get(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query().Get("key")

	w.Header().Add("Content-Type", "application/json")

	value, err := s.StringOps.Get(key)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(&types.StringOpsGetResponse{
		Key:   key,
		Value: value,
	})
}
