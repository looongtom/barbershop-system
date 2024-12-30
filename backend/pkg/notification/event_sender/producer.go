package main

import (
	"encoding/json"
	"fmt"
	"log"

	"DoAn/api"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	kafkaBroker = "localhost:9092"
	topic       = "booking"
)

// var kafkaBroker = os.Getenv("KAFKA_BROKER")
type Message struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	// Create a new Kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaBroker})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return
	}
	defer p.Close()
	// Define the message to be sent
	bookingRequest := api.BookingRequest{
		BarberId:      1,
		CustomerID:    2,
		SlotId:        4,
		Status:        "Pending",
		Price:         70000,
		ListServiceId: []int{1, 2, 3},
	}
	// Serialize the BookingRequest
	serializedBookingRequest, err := json.Marshal(bookingRequest)
	if err != nil {
		log.Fatalf("Failed to serialize booking request: %s\n", err)
	}

	if err != nil {
		fmt.Printf("Failed to serialize message: %s\n", err)
		return
	}
	// Produce the message to the Kafka topic
	err = produceMessage(p, topic, serializedBookingRequest)
	if err != nil {
		log.Fatalf("Failed to produce message: %s\n", err)
	}
	fmt.Println("Message produced successfully!")
}

// func serializeMessage(message Message) ([]byte, error) {
//
//	// Serialize the message struct to JSON
//
//	serialized, err := json.Marshal(message)
//
//	if err != nil {
//
//		return nil, fmt.Errorf("failed to serialize message: %w", err)
//
//	}
//
//	return serialized, nil
//
// }

func produceMessage(p *kafka.Producer, topic string, message []byte) error {
	// Create a new Kafka message to be produced
	kafkaMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}
	// Produce the Kafka message
	deliveryChan := make(chan kafka.Event)
	err := p.Produce(kafkaMessage, deliveryChan)
	if err != nil {

		return fmt.Errorf("failed to produce message: %w", err)

	}
	// Wait for delivery report or error
	e := <-deliveryChan
	m := e.(*kafka.Message)
	// Check for delivery errors
	if m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %s", m.TopicPartition.Error)
	}
	// Close the delivery channel
	close(deliveryChan)
	return nil
}
