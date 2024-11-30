package main

import (
	"encoding/json"
	"fmt"
	"log"
	logV "log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"

	"DoAn/api"
	"DoAn/database"
	"DoAn/entity"
)

const websocketURL = "ws://localhost:8080/trigger_hairfast" // Assuming WebSocket server is on the same machine

var (
	groupPreviewImageID     = "my_consumer_group"
	topicPreviewImage       = "preview_img"
	kafkaBrokerPreviewImage = "localhost:9092"

	hairfastResultTopic = "result_hairfast"
)

const (
	titleSuccessGenerate = "Generate successfully"
	titleFailGenerate    = "Generate failed"
	title                = "Choose to view more detail"
	typeNoti             = "hairfast"
)

func saveNotificationToMongo(result api.HairFastResult) error {
	titleNoti := titleSuccessGenerate
	if result.GeneratedImgCloud == "" {
		titleNoti = titleFailGenerate
	}
	encodeResult, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("Failed to encode booking: %s\n", err)
		encodeResult = nil
	}
	noti := entity.Notification{
		UserId:    1, // TODO: Hardcoded for now
		Title:     titleNoti,
		Message:   title,
		Type:      typeNoti,
		Timestamp: time.Now().Unix(),
		RawData:   encodeResult,
		IsRead:    false,
	}
	err = database.SaveNotification(noti)
	if err != nil {
		return err
	}
	fmt.Printf("Saved notification to MongoDB: %+v\n", noti)
	return nil
}

func sendToWebSocket(result api.HairFastResult) {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Printf("Failed to marshal HairFastResult to JSON: %v", err)
		return
	}
	// Create WebSocket client connection
	conn, _, err := websocket.DefaultDialer.Dial(websocketURL, http.Header{})
	if err != nil {
		log.Printf("Failed to connect to WebSocket server: %v", err)
		return
	}
	defer conn.Close()

	// Send the data
	err = conn.WriteMessage(websocket.TextMessage, resultJSON)
	if err != nil {
		log.Printf("Failed to send message to WebSocket server: %v", err)
	} else {
		log.Printf("Successfully sent message to WebSocket server: %s", resultJSON)
	}
}

func main() {

	err := godotenv.Load("notification.env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}

	kafkaBrokerPreviewImage = os.Getenv("KAFKA_BROKER")
	fmt.Printf("Kafka broker: %s\n", kafkaBrokerPreviewImage)

	kafkaBrokerServer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaBrokerPreviewImage})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return
	}
	defer kafkaBrokerServer.Close()

	// connGrpcBooking, err := grpc.Dial(os.Getenv("GRPC_BOOKING_SERVER"), grpc.WithInsecure(), grpc.WithBlock())
	// if err != nil {
	//	fmt.Printf("did not connect: %v", err)
	//	logV.Fatalf("Error getting env, %v", err)
	// }
	// client := pb.NewBookingServiceClient(connGrpcBooking)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBrokerPreviewImage,
		"group.id":          groupPreviewImageID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		return
	}
	defer c.Close()

	log.Println("Kafka listening on :9092 of topic " + hairfastResultTopic)

	// Subscribe to the Kafka topic
	err = c.SubscribeTopics([]string{hairfastResultTopic}, nil)
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
				var hairfastResult api.HairFastResult
				err := json.Unmarshal(e.Value, &hairfastResult)
				if err != nil {
					fmt.Printf("Failed to deserialize message: %s\n", err)
					continue
				}
				fmt.Printf("Received booking: %+v\n", hairfastResult)
				// call another api
				// CallAnotherAPI(previewImg)
				sendToWebSocket(hairfastResult)
				err = saveNotificationToMongo(hairfastResult)
				if err != nil {
					fmt.Printf("Failed to save notification to MongoDB: %s\n", err)
				}

			case kafka.Error:
				// Handle Kafka errors
				fmt.Printf("Error: %v\n", e)
			}
		}
	}
}
