package producer

import (
	"encoding/json"

	"kafka-example/models"

	"github.com/IBM/sarama" // Kafka client library
)

// Producer wraps the Sarama SyncProducer
type Producer struct {
	syncProducer sarama.SyncProducer // Underlying Kafka producer
}

// NewProducer initializes a new Kafka producer
func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()                     // Create a new Sarama configuration
	config.Producer.Return.Successes = true          // mandatory.. Ensure successes are returned
	config.Producer.RequiredAcks = sarama.WaitForAll // optional.. Wait for all replicas to acknowledge
	config.Producer.Retry.Max = 5                    // optional.. Set max retries for reliability
	config.Version = sarama.MaxVersion               // optional.. Specify Kafka version

	syncProducer, err := sarama.NewSyncProducer(brokers, config) // Initialize sync producer
	if err != nil {
		return nil, err // Return error if producer creation fails
	}

	return &Producer{syncProducer: syncProducer}, nil // Return wrapped producer
}

// SendMessage sends a message to the specified topic
func (p *Producer) SendMessage(topic string, msg models.Message) (int32, int64, error) {
	jsonValue, err := json.Marshal(msg) // Serialize Message struct to JSON
	if err != nil {
		return 0, 0, err // Return error if serialization fails
	}

	kafkaMsg := &sarama.ProducerMessage{ // Create a new producer message
		Topic: topic,                         // Set the target topic
		Key:   sarama.StringEncoder(msg.ID),  // Set the message key
		Value: sarama.ByteEncoder(jsonValue), // Set the JSON-encoded value
	}

	partition, offset, err := p.syncProducer.SendMessage(kafkaMsg) // Send the message synchronously
	return partition, offset, err                                  // Return partition, offset, and error
}

// Close shuts down the producer
func (p *Producer) Close() error {
	return p.syncProducer.Close() // Close the underlying producer
}
