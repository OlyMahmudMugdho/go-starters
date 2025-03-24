# Building a Kafka-Based RESTful API with Golang: A Step-by-Step Guide

In this comprehensive guide, we'll walk through the process of building a RESTful API in Golang that integrates with Apache Kafka. The API will allow users to send messages to a Kafka topic via an HTTP endpoint, and we'll also set up a Kafka consumer to process these messages asynchronously. We'll use the **Gorilla Mux** router for handling HTTP requests and the **IBM Sarama** library for Kafka integration. By the end of this guide, you'll have a modular, scalable, and production-ready application.

---

## Project Overview

Here's what we'll build:

- **Producer**: An HTTP endpoint (`POST /api/messages`) that accepts JSON payloads and sends them to a Kafka topic.
- **Consumer**: A background process that consumes messages from the Kafka topic and logs them.
- **Router**: Gorilla Mux for managing HTTP routes.
- **Error Handling**: Error management with logging.
- **Concurrency**: Efficient use of goroutines for asynchronous operations.

---

## Prerequisites

Before we begin, ensure you have the following set up:

- **Kafka**: A running Kafka broker (e.g., on `localhost:9092`). You can use Docker to start Kafka:
  ```bash
  docker run -p 9092:9092 apache/kafka:latest
  ```
- **Go**: Installed on your system with Go modules enabled.
- **Dependencies**: We'll install these as we go:
  ```bash
  go get github.com/IBM/sarama
  go get github.com/gorilla/mux
  ```

Initialize a new Go module for the project:
```bash
mkdir kafka-example
cd kafka-example
go mod init kafka-example
```

---

## Project Structure

To keep the code modular and maintainable, we'll organize it into the following structure:

```
kafka-example/
├── api/
│   └── server.go         # HTTP server and endpoint logic
├── consumer/
│   └── consumer.go       # Kafka consumer implementation
├── producer/
│   └── producer.go       # Kafka producer implementation
├── models/
│   └── message.go        # Data model for messages
├── go.mod                # Go module file
└── main.go               # Entry point to wire everything together
```

---

## Step 1: Define the Message Model

First, let's define the structure of the messages we'll send to Kafka.

### `models/message.go`

```go
package models

// Message defines the structure of data sent to Kafka
type Message struct {
    ID      string `json:"id"`      // Unique identifier for the message
    Content string `json:"content"` // Content of the message
}
```

### Explanation

- **Purpose**: This struct represents the data format for messages exchanged with Kafka.
- **Fields**:
  - `ID`: A unique identifier for each message.
  - `Content`: The actual message content.
- **JSON Tags**: The `json:"..."` tags ensure proper serialization/deserialization when handling HTTP requests and Kafka messages.

---

## Step 2: Implement the Kafka Producer

Next, we'll create a Kafka producer to send messages to a topic.

### `producer/producer.go`

```go
package producer

import (
    "encoding/json"
    "github.com/IBM/sarama"
    "kafka-example/models"
)

// Producer wraps the Sarama SyncProducer
type Producer struct {
    syncProducer sarama.SyncProducer
}

// NewProducer initializes a new Kafka producer
func NewProducer(brokers []string) (*Producer, error) {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true        // Wait for successful delivery
    config.Producer.RequiredAcks = sarama.WaitForAll // Require all replicas to acknowledge
    config.Producer.Retry.Max = 5                  // Retry up to 5 times on failure
    config.Version = sarama.MaxVersion             // Use the latest Kafka version

    syncProducer, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        return nil, err
    }
    return &Producer{syncProducer: syncProducer}, nil
}

// SendMessage sends a message to the specified topic
func (p *Producer) SendMessage(topic string, msg models.Message) (int32, int64, error) {
    jsonValue, err := json.Marshal(msg)
    if err != nil {
        return 0, 0, err
    }
    kafkaMsg := &sarama.ProducerMessage{
        Topic: topic,
        Key:   sarama.StringEncoder(msg.ID),  // Use ID as the message key
        Value: sarama.ByteEncoder(jsonValue), // JSON-encoded message
    }
    partition, offset, err := p.syncProducer.SendMessage(kafkaMsg)
    return partition, offset, err
}

// Close shuts down the producer
func (p *Producer) Close() error {
    return p.syncProducer.Close()
}
```

### Explanation

- **Producer Struct**: Wraps a `sarama.SyncProducer` for synchronous message sending.
- **NewProducer**:
  - **Configuration**:
    - `Return.Successes`: Ensures we get confirmation of successful delivery.
    - `RequiredAcks`: Waits for all Kafka replicas to acknowledge the message for reliability.
    - `Retry.Max`: Retries up to 5 times if sending fails.
    - `Version`: Uses the latest Kafka protocol version for compatibility.
  - **Initialization**: Creates a new synchronous producer connected to the specified Kafka brokers.
- **SendMessage**:
  - Serializes the `Message` struct to JSON.
  - Constructs a `ProducerMessage` with the topic, key (message ID), and value (JSON bytes).
  - Sends the message and returns the partition and offset where it was stored.
- **Close**: Properly shuts down the producer to free resources.

---

## Step 3: Implement the Kafka Consumer

Now, let's build a consumer to process messages from the Kafka topic.

### `consumer/consumer.go`

```go
package consumer

import (
    "context"
    "encoding/json"
    "log"
    "time"
    "github.com/IBM/sarama"
    "kafka-example/models"
)

// ConsumerGroupHandler implements the Sarama ConsumerGroupHandler interface
type ConsumerGroupHandler struct{}

// Setup is called before consuming starts
func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
    return nil
}

// Cleanup is called after consuming stops
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
    return nil
}

// ConsumeClaim processes messages from a partition
func (ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        var message models.Message
        err := json.Unmarshal(msg.Value, &message)
        if err != nil {
            log.Printf("Failed to deserialize message: %v", err)
            continue
        }
        log.Printf("Consumed message: id=%s, content=%s, topic=%s, partition=%d, offset=%d",
            message.ID, message.Content, msg.Topic, msg.Partition, msg.Offset)
        session.MarkMessage(msg, "") // Mark message as processed
    }
    return nil
}

// StartConsumer launches the Kafka consumer in a goroutine
func StartConsumer(brokers []string, group, topic string) {
    config := sarama.NewConfig()
    config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
    config.Consumer.Offsets.Initial = sarama.OffsetOldest // Start from the oldest messages
    config.Version = sarama.MaxVersion

    client, err := sarama.NewConsumerGroup(brokers, group, config)
    if err != nil {
        log.Fatalf("Error creating consumer group client: %v", err)
    }
    defer client.Close()

    handler := ConsumerGroupHandler{}
    ctx := context.Background()

    go func() {
        for {
            err := client.Consume(ctx, []string{topic}, handler)
            if err != nil {
                log.Printf("Error from consumer: %v", err)
            }
            time.Sleep(1 * time.Second) // Brief pause before retrying
        }
    }()
}
```

### Explanation

- **ConsumerGroupHandler**: Implements the `sarama.ConsumerGroupHandler` interface:
  - **Setup**: Runs before consuming begins (no-op here).
  - **Cleanup**: Runs after consuming stops (no-op here).
  - **ConsumeClaim**: Processes messages from a partition:
    - Deserializes the JSON message into a `Message` struct.
    - Logs the message details (ID, content, topic, partition, offset).
    - Marks the message as processed in the session.
- **StartConsumer**:
  - **Configuration**:
    - `Rebalance.Strategy`: Uses round-robin to distribute partitions among consumers.
    - `Offsets.Initial`: Starts consuming from the oldest available messages.
  - **Initialization**: Creates a consumer group client.
  - **Goroutine**: Runs the consumer in a background loop, retrying on errors with a 1-second delay.

---

## Step 4: Set Up the HTTP Server

Let's create an HTTP server to expose an endpoint for producing messages.

### `api/server.go`

```go
package api

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
    "github.com/gorilla/mux"
    "kafka-example/models"
    "kafka-example/producer"
)

// Server holds the application state and dependencies
type Server struct {
    Producer *producer.Producer
    Topic    string
}

// NewServer creates a new server instance
func NewServer(p *producer.Producer, topic string) *Server {
    return &Server{Producer: p, Topic: topic}
}

// ProduceMessageHandler handles POST requests to send messages to Kafka
func (s *Server) ProduceMessageHandler(w http.ResponseWriter, r *http.Request) {
    var msg models.Message
    if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }
    partition, offset, err := s.Producer.SendMessage(s.Topic, msg)
    if err != nil {
        log.Printf("Failed to send message to Kafka: %v", err)
        http.Error(w, "Failed to produce message", http.StatusInternalServerError)
        return
    }
    log.Printf("Message sent to partition %d at offset %d", partition, offset)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"status": "message produced"})
}

// StartServer starts the HTTP server
func StartServer(p *producer.Producer, topic string) {
    server := NewServer(p, topic)
    router := mux.NewRouter()
    router.HandleFunc("/api/messages", server.ProduceMessageHandler).Methods("POST")
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      router,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
    log.Println("Starting server on :8080")
    if err := srv.ListenAndServe(); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

### Explanation

- **Server Struct**: Holds the Kafka producer and topic name for dependency injection.
- **ProduceMessageHandler**:
  - Decodes the JSON request body into a `Message` struct.
  - Sends the message to Kafka using the producer.
  - Returns a success response (`201 Created`) or an error if something fails.
- **StartServer**:
  - Sets up a Gorilla Mux router with the `/api/messages` endpoint.
  - Configures an HTTP server with timeouts and starts it on port `8080`.

---

## Step 5: Wire Everything Together

Finally, let's tie everything together in the main entry point.

### `main.go`

```go
package main

import (
    "log"
    "kafka-example/api"
    "kafka-example/consumer"
    "kafka-example/producer"
)

func main() {
    brokers := []string{"localhost:9092"}
    topic := "messages"
    consumerGroup := "message-consumer-group"

    // Initialize Kafka producer
    p, err := producer.NewProducer(brokers)
    if err != nil {
        log.Fatalf("Failed to initialize producer: %v", err)
    }
    defer p.Close()

    // Start Kafka consumer in the background
    consumer.StartConsumer(brokers, consumerGroup, topic)

    // Start HTTP server
    api.StartServer(p, topic)
}
```

### Explanation

- **Configuration**: Defines Kafka brokers, topic, and consumer group name.
- **Producer**: Initializes the Kafka producer and ensures it’s closed when the program exits.
- **Consumer**: Starts the consumer in a goroutine to run concurrently.
- **Server**: Launches the HTTP server, which blocks the main thread.

---

## Testing the Application

Let’s test the application to ensure everything works as expected.

1. **Start Kafka**: Ensure your Kafka broker is running (e.g., via Docker).
2. **Run the Application**:
   ```bash
   go run main.go
   ```
   You should see: `Starting server on :8080`.
3. **Send a Message**:
   Use `curl` to send a message to the API:
   ```bash
   curl -X POST -H "Content-Type: application/json" -d '{"id":"123","content":"Hello, Kafka!"}' http://localhost:8080/api/messages
   ```
   Expected response:
   ```json
   {"status":"message produced"}
   ```
4. **Check Logs**:
   - **Producer**: `Message sent to partition X at offset Y`
   - **Consumer**: `Consumed message: id=123, content=Hello, Kafka!, topic=messages, partition=X, offset=Y`

---

## Conclusion

Congratulations! You've built a Kafka-based RESTful API in Golang. Here's what we accomplished:

- **Modularity**: Separated concerns into `models`, `producer`, `consumer`, and `api` packages.
- **Reliability**: Configured the producer and consumer for robust operation with error handling and retries.
- **Concurrency**: Used goroutines to handle the consumer asynchronously.
- **Scalability**: The design supports adding more endpoints or consumers as needed.

This application is a solid foundation for production use. You can extend it by adding more endpoints, enhancing error handling, or integrating with a database. Happy coding!