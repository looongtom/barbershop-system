package main

import (
	"DoAn/pkg/booking/api"
	kafka2 "DoAn/pkg/booking/kafka"
	"DoAn/pkg/booking/pb"
	"encoding/json"
	"fmt"
	logV "log"
	"os"
	"os/signal"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

const (
	groupID     = "booking-group"
	topic       = "booking"
	replyTopic  = "reply"
	kafkaBroker = "localhost:9092"
)

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

func main() {

	err := godotenv.Load("criteria.env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}

	kafkaBrokerServer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaBroker})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return
	}
	defer kafkaBrokerServer.Close()

	//connGrpcBooking, err := grpc.Dial(os.Getenv("GRPC_BOOKING_SERVER"), grpc.WithInsecure(), grpc.WithBlock())
	//if err != nil {
	//	fmt.Printf("did not connect: %v", err)
	//	logV.Fatalf("Error getting env, %v", err)
	//}
	//client := pb.NewBookingServiceClient(connGrpcBooking)

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
				var booking api.KafkaBookingRequest
				err := json.Unmarshal(e.Value, &booking)
				if err != nil {
					fmt.Printf("Failed to deserialize message: %s\n", err)
					continue
				}
				fmt.Printf("Received booking: %+v\n", booking)

				//// Convert []int32 to []int
				//listServiceId := make([]int32, len(booking.ListServiceId))
				//for i, v := range booking.ListServiceId {
				//	listServiceId[i] = int32(v)
				//}
				//createBooking, err := client.CreateBooking(context.Background(), &pb.BookingRequest{
				//	CustomerId:    int32(booking.CustomerID),
				//	BarberId:      int32(booking.BarberId),
				//	Status:        booking.Status,
				//	Price:         booking.Price,
				//	SlotId:        int32(booking.SlotId),
				//	ListServiceId: listServiceId,
				//})
				//if err != nil {
				//	fmt.Printf("error while creating booking: %v\n", err)
				//	sendKafkaResponse(kafkaBrokerServer, &pb.Booking{}, replyTopic, booking.UUID)
				//	continue
				//}
				//fmt.Printf("Created booking successfully: %+v\n", createBooking)
				//sendKafkaResponse(kafkaBrokerServer, createBooking, replyTopic, booking.UUID)
			case kafka.Error:
				// Handle Kafka errors
				fmt.Printf("Error: %v\n", e)
			}
		}
	}
}
