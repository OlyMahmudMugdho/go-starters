package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"kafka-example/models"
	"kafka-example/producer"

	"github.com/gorilla/mux"
)

// Server holds the application state and dependencies
type Server struct {
	Producer *producer.Producer // Kafka producer instance
	Topic    string             // Kafka topic name
}

// NewServer creates a new server instance
func NewServer(p *producer.Producer, topic string) *Server {
	return &Server{Producer: p, Topic: topic} // Initialize server with producer and topic
}

// ProduceMessageHandler handles POST requests to send messages to Kafka
func (s *Server) ProduceMessageHandler(w http.ResponseWriter, r *http.Request) {
	var msg models.Message                                       // Define message variable
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil { // Decode JSON payload
		http.Error(w, "Invalid request payload", http.StatusBadRequest) // Return error on invalid payload
		return
	}

	partition, offset, err := s.Producer.SendMessage(s.Topic, msg) // Send message to Kafka
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)                     // Log error
		http.Error(w, "Failed to produce message", http.StatusInternalServerError) // Return error response
		return
	}

	log.Printf("Message sent to partition %d at offset %d", partition, offset) // Log success
	w.WriteHeader(http.StatusCreated)                                          // Set status to 201
	json.NewEncoder(w).Encode(map[string]string{"status": "message produced"}) // Send success response
}

// StartServer starts the HTTP server
func StartServer(p *producer.Producer, topic string) {
	server := NewServer(p, topic)                                                    // Create server instance
	router := mux.NewRouter()                                                        // Initialize Gorilla Mux router
	router.HandleFunc("/api/messages", server.ProduceMessageHandler).Methods("POST") // Register endpoint

	srv := &http.Server{ // Configure HTTP server
		Addr:         ":8080",          // Listen on port 8080
		Handler:      router,           // Use router
		ReadTimeout:  10 * time.Second, // Set read timeout
		WriteTimeout: 10 * time.Second, // Set write timeout
	}

	log.Println("Starting server on :8080")      // Log server start
	if err := srv.ListenAndServe(); err != nil { // Start server and handle errors
		log.Fatalf("Server failed: %v", err) // Fatal error if server fails
	}
}
