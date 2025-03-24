package models

// Message defines the structure of data sent to Kafka
type Message struct {
	ID      string `json:"id"`      // Unique identifier for the message
	Content string `json:"content"` // Content of the message
}
