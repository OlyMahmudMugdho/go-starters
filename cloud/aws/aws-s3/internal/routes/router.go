package routes

import (
	"aws-s3/internal/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter() http.Handler {
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/s3/create", handler.CreateBucketHandler).Methods("POST")
	r.HandleFunc("/s3/list", handler.ListBucketsHandler).Methods("GET")

	r.HandleFunc("/s3/delete/{name}", handler.DeleteBucketHandler).Methods("DELETE")
	r.HandleFunc("/s3/update", handler.UpdateBucketHandler).Methods("PUT")

	return r
}
