package handlers

import (
	"encoding/json"
	"go-redis-practice/types"
	"go-redis-practice/utils"
	"net/http"
	"strconv"
)

type ListHandler struct {
	ListOps *utils.ListOps
}

func NewListHandler() *ListHandler {
	ListOps := utils.NewListOps()
	return &ListHandler{
		ListOps: ListOps,
	}
}

func (l *ListHandler) LPush(w http.ResponseWriter, r *http.Request) {
	var body types.List

	json.NewDecoder(r.Body).Decode(&body)

	w.Header().Add("Content-Type", "application/json")

	response, err := l.ListOps.LPush(body.Key, body.Value)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatInt(response, 10)))
}

func (l *ListHandler) LRange(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	start := r.URL.Query().Get("start")
	stop := r.URL.Query().Get("stop")

	w.Header().Add("Content-Type", "application/json")

	startInt, _ := strconv.Atoi(start)
	stopInt, _ := strconv.Atoi(stop)

	response, err := l.ListOps.LRange(key, int64(startInt), int64(stopInt))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatInt(int64(len(response)), 10)))

}

func (l *ListHandler) LLen(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	w.Header().Add("Content-Type", "application/json")

	response, err := l.ListOps.LLen(key)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatInt(response, 10)))
}

func (l *ListHandler) LPop(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query().Get("key")
	w.Header().Add("Content-Type", "application/json")

	response, err := l.ListOps.LPop(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func (l *ListHandler) RPush(w http.ResponseWriter, r *http.Request) {
	var body types.List

	json.NewDecoder(r.Body).Decode(&body)

	w.Header().Add("Content-Type", "application/json")

	response, err := l.ListOps.RPush(body.Key, body.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatInt(response, 10)))
}

func (l *ListHandler) RPop(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	w.Header().Add("Content-Type", "application/json")

	response, err := l.ListOps.RPop(key)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func (l *ListHandler) LIndex(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	index := r.URL.Query().Get("index")

	indexInt, _ := strconv.Atoi(index)

	w.Header().Add("Content-Type", "application/json")

	response, err := l.ListOps.LIndex(key, int64(indexInt))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
