package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	logV "log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"

	"DoAn/api"
	"DoAn/database"
	"DoAn/entity"
	kafka2 "DoAn/kafka"
	"DoAn/pb"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

var (
	groupID      = "booking-group"
	topic        = "booking"
	kafkaBroker  = "localhost:9092"
	websocketURL = "ws://localhost:8080/trigger_booking"
)

const (
	titleSuccessBooking = "Booking successfully"
	titleFailBooking    = "Booking failed"
	title               = "Choose to view more detail"
	typeNoti            = "booking"
)

func saveNotificationToMongo(bookingResp api.BookingResponse) error {
	var titleNoti string
	switch bookingResp.Status {
	case "Booked":
		titleNoti = titleSuccessBooking
	default:
		titleNoti = titleFailBooking
	}
	encodeBooking, err := json.Marshal(bookingResp)
	if err != nil {
		fmt.Printf("Failed to encode booking: %s\n", err)
		encodeBooking = nil
	}

	noti := entity.Notification{
		UserId:    bookingResp.CustomerID,
		Title:     titleNoti,
		Message:   title,
		Type:      typeNoti,
		Timestamp: bookingResp.CreatedAt,
		RawData:   encodeBooking,
		IsRead:    false,
	}
	// Save notification to MongoDB
	err = database.SaveNotification(noti)
	if err != nil {
		return err
	}
	fmt.Printf("Saved notification to MongoDB: %+v\n", noti)

	noti2 := entity.Notification{
		UserId:    bookingResp.BarberId,
		Title:     titleNoti,
		Message:   title,
		Type:      typeNoti,
		Timestamp: bookingResp.CreatedAt,
		RawData:   encodeBooking,
		IsRead:    false,
	}
	// Save notification to MongoDB
	err = database.SaveNotification(noti2)
	if err != nil {
		return err
	}
	fmt.Printf("Saved notification to MongoDB: %+v\n", noti2)
	return nil
}

func sendKafkaResponse(kaf *kafka.Producer, createBooking *pb.Booking, topic, uuid string) {

	bookingResponse := api.BookingResponse{
		ID:         int(createBooking.Id),
		CustomerID: int(createBooking.CustomerId),
		BarberId:   int(createBooking.BarberId),
		Status:     createBooking.Status,
		Price:      createBooking.Price,
		SlotId:     int(createBooking.SlotId),
		CreatedAt:  int64(createBooking.CreatedAt),
		UpdatedAt:  int64(createBooking.UpdatedAt),
	}
	// Serialize the BookingRequest
	serializedBookingRequest, err := json.Marshal(bookingResponse)
	if err != nil {
		logV.Fatalf("Failed to serialize booking request: %s\n", err)
	}

	if err != nil {
		fmt.Printf("Failed to serialize message: %s\n", err)
		return
	}
	// Produce the message to the Kafka topic
	err = kafka2.ProduceMessage(kaf, topic, serializedBookingRequest)
	if err != nil {
		logV.Fatalf("Failed to produce message: %s\n", err)
	}
	fmt.Println("Message produced successfully!")
}

func sendToWebSocket(result api.BookingResponse) {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Printf("Failed to marshal Booking to JSON: %v", err)
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

func convertListInt32ToInt(listInt32 []int32) []int {
	listInt := make([]int, len(listInt32))
	for i, v := range listInt32 {
		listInt[i] = int(v)
	}
	return listInt
}

func convertListIntToInt32(listInt []int) []int32 {
	listInt32 := make([]int32, len(listInt))
	for i, v := range listInt {
		listInt32[i] = int32(v)
	}
	return listInt32
}

func convertInt32ToString(listInt32 []int32) []string {
	listString := make([]string, len(listInt32))
	for i, v := range listInt32 {
		listString[i] = strconv.Itoa(int(v))
	}
	return listString
}

func main() {

	err := godotenv.Load("notification.env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}

	kafkaBroker = os.Getenv("KAFKA_BROKER")

	kafkaBrokerServer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaBroker})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return
	}
	defer kafkaBrokerServer.Close()

	connGrpcBooking, err := grpc.Dial(os.Getenv("GRPC_BOOKING_SERVER"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("did not connect: %v", err)
		logV.Fatalf("Error getting env, %v", err)
	}
	client := pb.NewBookingServiceClient(connGrpcBooking)

	// connGrpcBooking, err := grpc.Dial(os.Getenv("GRPC_BOOKING_SERVER"), grpc.WithInsecure(), grpc.WithBlock())
	// if err != nil {
	//	fmt.Printf("did not connect: %v", err)
	//	logV.Fatalf("Error getting env, %v", err)
	// }
	// client := pb.NewBookingServiceClient(connGrpcBooking)

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
	fmt.Printf("Consuming messages from topic: %s\n", topic)
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
				var booking api.KafkaBookingRequest
				err := json.Unmarshal(e.Value, &booking)
				if err != nil {
					fmt.Printf("Failed to deserialize message: %s\n", err)
					continue
				}
				fmt.Printf("Received booking: %+v\n", booking)

				listServiceId32 := make([]int32, len(booking.ListServiceId))
				for i, v := range booking.ListServiceId {
					listServiceId32[i] = int32(v)
				}
				resp, err := client.CreateBooking(context.Background(), &pb.BookingRequest{
					CustomerId:    int32(booking.CustomerID),
					BarberId:      int32(booking.BarberId),
					Status:        "Booked",
					Price:         booking.Price,
					SlotId:        int32(booking.SlotId),
					ListServiceId: listServiceId32,
					BookedDate:    booking.BookedDate,
				})
				if err != nil {
					fmt.Printf("error while creating booking: %v\n", err)
					continue
				}
				fmt.Printf("Created booking successfully: %+v\n", resp)
				successBookingResp := api.BookingResponse{
					ID:         int(resp.Id),
					CustomerID: int(resp.CustomerId),
					BarberId:   int(resp.BarberId),
					ResultId:   int(resp.ResultId),
					Status:     resp.Status,
					Price:      resp.Price,
					SlotId:     int(resp.SlotId),
					FeedBackId: int(resp.FeedbackId),
					CreatedAt:  int64(resp.CreatedAt),
					UpdatedAt:  int64(resp.UpdatedAt),
					// ListServices: convertInt32ToString(resp.ListServiceId),
				}
				sendToWebSocket(successBookingResp)
				err = saveNotificationToMongo(successBookingResp)
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
