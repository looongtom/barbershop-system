package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	logV "log"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"

	"DoAn/pkg/previewimage/api"
)

const (
	groupPreviewImageID     = "preview_img-group"
	topicPreviewImage       = "preview_img"
	kafkaBrokerPreviewImage = "localhost:9092"
)

func main() {

	err := godotenv.Load("criteria.env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}

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

	// Subscribe to the Kafka topic
	err = c.SubscribeTopics([]string{topicPreviewImage}, nil)
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
				var previewImg api.KafkaPreviewImageRequest
				err := json.Unmarshal(e.Value, &previewImg)
				if err != nil {
					fmt.Printf("Failed to deserialize message: %s\n", err)
					continue
				}
				fmt.Printf("Received booking: %+v\n", previewImg)
				// call another api
				CallAnotherAPI(previewImg)
			case kafka.Error:
				// Handle Kafka errors
				fmt.Printf("Error: %v\n", e)
			}
		}
	}
}

func CallAnotherAPI(previewImg api.KafkaPreviewImageRequest) {
	url := "http://192.168.1.3:5000/upload-images"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	var wg sync.WaitGroup

	downloadAndCreateFormFile := func(fieldName, imageURL string) {
		defer wg.Done()
		resp, err := http.Get(imageURL)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("failed to download image: %s\n", resp.Status)
			return
		}
		part, err := writer.CreateFormFile(fieldName, fieldName)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = io.Copy(part, resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	wg.Add(3)
	go downloadAndCreateFormFile("selfImg", previewImg.SelfImg)
	go downloadAndCreateFormFile("shapeImg", previewImg.ShapeImg)
	go downloadAndCreateFormFile("colorImg", previewImg.ColorImg)

	wg.Wait()

	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
