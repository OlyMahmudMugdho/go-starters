package main

import (
	"log"
	"net/http"
	"simple-websocket/server"
)

func main() {
	// Create a new WebSocket server
	s := server.NewWebSocketServer()

	// Start the WebSocket server's event loop in a goroutine
	go s.Run()

	// Define the WebSocket endpoint
	http.HandleFunc("/ws", s.HandleWebSocket)

	// Serve static files (HTML, CSS, JS)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Start the server on port 8080
	log.Println("Starting WebSocket server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
