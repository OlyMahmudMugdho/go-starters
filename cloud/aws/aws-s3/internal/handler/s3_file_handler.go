package handler

import (
	"aws-s3/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.FormValue("bucket")
	key := r.FormValue("key")
	file, _, err := r.FormFile("file")

	if err != nil || bucket == "" || key == "" {
		http.Error(w, "Invalid form input", http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = service.UploadFile(context.Background(), bucket, key, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "File uploaded successfully")
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucketName := vars["bucket"]
	key := vars["key"]

	fileStream, err := service.DownloadFile(context.Background(), bucketName, key)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound) // Use 404 for NoSuchKey errors
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to download file: %s", err.Error()),
		})
		return
	}
	defer fileStream.Close()

	w.Header().Set("Content-Disposition", `attachment; filename="`+key+`"`)
	w.Header().Set("Content-Type", "application/octet-stream")

	_, err = io.Copy(w, fileStream)
	if err != nil {
		// Log the error if needed, but the response is already partially sent
		log.Printf("Error copying file stream to response: %v", err)
	}
}

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	bucket := mux.Vars(r)["bucket"]
	key := mux.Vars(r)["key"]

	err := service.DeleteFile(context.Background(), bucket, key)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to delete file: %v", err),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("File '%s' deleted from bucket '%s'", key, bucket),
	})
}
