package main

import (
	"DoAn/pkg/previewimage/api"
	"encoding/json"
	"fmt"
	logV "log"
	"os"
	"os/signal"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

const (
	groupID     = "preview-img-group"
	topic       = "preview-img"
	kafkaBroker = "localhost:9092"
)

func main() {

	err := godotenv.Load("account.env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}

	kafkaBrokerServer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaBroker})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return
	}
	defer kafkaBrokerServer.Close()

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBroker,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		return
	}
	defer c.Close()

	// Subscribe to the Kafka topic
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topic: %s\n", err)
		return
	}

	// Setup a channel to handle OS signals for graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	// Start consuming messages
	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Received signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				// Process the consumed message
				var imgs api.UpdateImageRequest
				err := json.Unmarshal(e.Value, &imgs)
				if err != nil {
					fmt.Printf("Failed to deserialize message: %s\n", err)
					continue
				}
				fmt.Printf("Received imgs: %+v\n", imgs)

			case kafka.Error:
				// Handle Kafka errors
				fmt.Printf("Error: %v\n", e)
			}
		}
	}
}
