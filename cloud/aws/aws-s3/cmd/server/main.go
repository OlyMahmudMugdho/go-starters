package main

import (
	"aws-s3/internal/aws"
	"aws-s3/internal/routes"
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	ctx := context.Background()

	// Load default AWS config
	cfg, err := config.LoadDefaultConfig(ctx)
	cfg.Region = "eu-north-1"
	if err != nil {
		log.Fatalf("Unable to load SDK config: %v", err)
	}

	// Init S3 client
	aws.InitS3Client(cfg)

	// Setup and start server
	r := routes.SetupRouter()
	log.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
