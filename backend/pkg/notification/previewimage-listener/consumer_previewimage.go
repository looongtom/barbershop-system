package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
	"io"
	"log"
	logV "log"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"DoAn/api"
)

const (
	groupPreviewImageID     = "preview_img-group"
	topicPreviewImage       = "preview_img"
	kafkaBrokerPreviewImage = "localhost:9092"

	hairfastTopic = "hairfast"
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
				fmt.Printf("Received previewImg: %+v\n", previewImg)
				// call another api
				//CallAnotherAPI(previewImg)

				// Produce the message to the Kafka topic
				serializedPreviewImage, err := json.Marshal(previewImg)
				if err != nil {
					log.Fatalf("Failed to serialize booking request: %s\n", err)
					return
				}

				err = produceMessage(kafkaBrokerServer, hairfastTopic, serializedPreviewImage)
				if err != nil {
					log.Fatalf("Failed to produce message: %s\n", err)
				}
				fmt.Println("Message produced successfully!")

			case kafka.Error:
				// Handle Kafka errors
				fmt.Printf("Error: %v\n", e)
			}
		}
	}
}

func produceMessage(p *kafka.Producer, topic string, message []byte) error {
	deliveryChan := make(chan kafka.Event)
	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, deliveryChan)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
	}
	fmt.Printf("Produced message to topic %s: %s\n", *m.TopicPartition.Topic, string(m.Value))
	return nil

}

func downloadAndCreateFormFile(writer *multipart.Writer, fieldName, imageURL string) error {
	resp, err := http.Get(imageURL)
	if err != nil {
		fmt.Println("1")
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("failed to download image: %s\n", resp.Status)
		return err
	}
	part, err := writer.CreateFormFile(fieldName, fieldName)
	if err != nil {
		fmt.Println("2")
		fmt.Println(err)
		return err
	}
	_, err = io.Copy(part, resp.Body)
	if err != nil {
		fmt.Println("3")
		fmt.Println(err)
		return err
	}
	return nil
}

func CallAnotherAPI(previewImg api.KafkaPreviewImageRequest) {
	url := "http://localhost:5000/upload-images"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	wg.Add(3)

	go func() {
		defer wg.Done()
		if err := downloadAndCreateFormFile(writer, "selfImg", previewImg.SelfImg); err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := downloadAndCreateFormFile(writer, "shapeImg", previewImg.ShapeImg); err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := downloadAndCreateFormFile(writer, "colorImg", previewImg.ColorImg); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		for err := range errChan {
			fmt.Println(err)
		}
		return
	}

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
