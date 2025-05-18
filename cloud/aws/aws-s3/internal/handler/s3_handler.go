package handler

import (
	"aws-s3/internal/service"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type CreateBucketRequest struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

type UpdateBucketRequest struct {
	OldName   string `json:"old_name"`
	NewName   string `json:"new_name"`
	NewRegion string `json:"new_region"`
}

func CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateBucketRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Name == "" || req.Region == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = service.CreateBucket(context.Background(), req.Name, req.Region)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Bucket created successfully",
	})
}

func ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	buckets, err := service.ListBuckets(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(buckets)
}

func DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}

	err := service.DeleteBucket(context.Background(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Bucket deleted successfully",
	})
}

func UpdateBucketHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateBucketRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.OldName == "" || req.NewName == "" || req.NewRegion == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = service.UpdateBucket(context.Background(), req.OldName, req.NewName, req.NewRegion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Bucket updated successfully",
	})
}
