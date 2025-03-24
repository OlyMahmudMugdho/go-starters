package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"kafka-example/models"

	"github.com/IBM/sarama" // Kafka client library
)

// ConsumerGroupHandler implements the Sarama ConsumerGroupHandler interface
type ConsumerGroupHandler struct{}

// Setup is called before consuming starts
func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil // No setup needed
}

// Cleanup is called after consuming stops
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil // No cleanup needed
}

// ConsumeClaim processes messages from a partition
func (ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() { // Iterate over messages in the claim
		var message models.Message                 // Define struct to hold deserialized message
		err := json.Unmarshal(msg.Value, &message) // Deserialize JSON value into Message struct
		if err != nil {
			log.Printf("Failed to deserialize message: %v", err) // Log deserialization error
			continue                                             // Skip to next message on error
		}

		log.Printf("Consumed message: id=%s, content=%s, topic=%s, partition=%d, offset=%d",
			message.ID, message.Content, msg.Topic, msg.Partition, msg.Offset) // Log struct fields
		session.MarkMessage(msg, "") // Mark message as processed
	}
	return nil // Return nil on success
}

// StartConsumer launches the Kafka consumer in a goroutine
func StartConsumer(brokers []string, group, topic string) {
	config := sarama.NewConfig()                                                     // Create a new Sarama configuration
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin() // optional.. Use updated round-robin strategy
	config.Consumer.Offsets.Initial = sarama.OffsetOldest                            // optional.. Start from oldest messages
	config.Version = sarama.MaxVersion                                               // optional.. Specify Kafka version

	client, err := sarama.NewConsumerGroup(brokers, group, config) // Initialize consumer group
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v", err) // Fatal error if client creation fails
	}

	handler := ConsumerGroupHandler{} // Create handler instance
	ctx := context.Background()       // Use background context

	go func() { // Start consumer in a goroutine
		for {
			err := client.Consume(ctx, []string{topic}, handler) // Consume messages from topic
			if err != nil {
				log.Printf("Error from consumer: %v", err) // Log consumption errors
			}
			time.Sleep(1 * time.Second) // Backoff before retrying
		}
	}()
}
