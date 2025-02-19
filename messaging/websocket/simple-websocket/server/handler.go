package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (you can customize this for security)
	},
}

// HandleWebSocket upgrades the HTTP connection to a WebSocket connection.
func (s *WebSocketServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Register the new client
	s.Register <- conn

	// Start listening for messages from the client
	go s.readMessages(conn)
}

// readMessages reads messages from a WebSocket connection.
func (s *WebSocketServer) readMessages(conn *websocket.Conn) {
	defer func() {
		s.Unregister <- conn
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Broadcast the received message to all clients except the sender
		go func(msg []byte) {
			s.mu.Lock()
			defer s.mu.Unlock()

			for client := range s.Clients {
				if client != conn { // Don't send the message back to the sender
					err := client.WriteMessage(websocket.TextMessage, msg)
					if err != nil {
						log.Printf("Error sending message: %v", err)
						client.Close()
						delete(s.Clients, client)
					}
				}
			}
		}(message)
	}
}
