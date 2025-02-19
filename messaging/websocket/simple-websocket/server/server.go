package server

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketServer holds the state of the WebSocket server.
type WebSocketServer struct {
	Clients    map[*websocket.Conn]bool
	Broadcast  chan []byte
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
	mu         sync.Mutex
}

// NewWebSocketServer initializes a new WebSocket server.
func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		Clients:    make(map[*websocket.Conn]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}
}

// Run starts the WebSocket server's event loop.
func (s *WebSocketServer) Run() {
	for {
		select {
		case conn := <-s.Register:
			s.mu.Lock()
			s.Clients[conn] = true
			s.mu.Unlock()
			log.Println("New client connected")

		case conn := <-s.Unregister:
			s.mu.Lock()
			if _, ok := s.Clients[conn]; ok {
				delete(s.Clients, conn)
				conn.Close()
				log.Println("Client disconnected")
			}
			s.mu.Unlock()

		case message := <-s.Broadcast:
			s.mu.Lock()
			for conn := range s.Clients {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error sending message: %v", err)
					conn.Close()
					delete(s.Clients, conn)
				}
			}
			s.mu.Unlock()
		}
	}
}
