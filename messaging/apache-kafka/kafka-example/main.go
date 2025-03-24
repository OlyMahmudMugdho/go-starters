package main

import (
	"log"

	"kafka-example/api"
	"kafka-example/consumer"
	"kafka-example/producer"
)

func main() {
	brokers := []string{"localhost:9092"}     // Kafka broker list
	topic := "messages"                       // Kafka topic name
	consumerGroup := "message-consumer-group" // Consumer group name

	p, err := producer.NewProducer(brokers) // Initialize Kafka producer
	if err != nil {
		log.Fatalf("Failed to initialize producer: %v", err) // Fatal error if producer fails
	}
	defer p.Close() // Ensure producer is closed on exit

	consumer.StartConsumer(brokers, consumerGroup, topic) // Start Kafka consumer

	api.StartServer(p, topic) // Start HTTP server
}
