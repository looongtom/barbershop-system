package service

import (
	"DoAn/pkg/account/pb"
	"DoAn/pkg/booking"
	"DoAn/pkg/booking/api"
	kafka2 "DoAn/pkg/booking/kafka"
	pbService "DoAn/pkg/servicing/pb"
	pbTimeslot "DoAn/pkg/timeslot/pb"
	"encoding/json"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"os"
	"os/signal"

	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
)

const (
	groupID     = "reply-group"
	replyTopic  = "reply"
	kafkaBroker = "localhost:9092"
)

type BookingStruct struct {
	repository   booking.BookingRepository
	connAccount  *grpc.ClientConn
	connTimeslot *grpc.ClientConn
	connService  *grpc.ClientConn
	logger       log.Logger
	kafka        *kafka.Producer
}

func NewService(repo booking.BookingRepository, logger log.Logger,
	conn *grpc.ClientConn, conn2 *grpc.ClientConn, conn3 *grpc.ClientConn,
	kafka *kafka.Producer) booking.BookingService {
	return &BookingStruct{
		repository:   repo,
		logger:       logger,
		connAccount:  conn,
		connTimeslot: conn2,
		connService:  conn3,
		kafka:        kafka,
	}
}

func (b BookingStruct) FindBookingByUserOrBarber(ctx context.Context, findReq api.FindBookingRequest) (interface{}, error) {
	resp, err := b.repository.FindBookingByUserOrBarber(ctx, findReq)
	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}
	return resp, nil
}

func (b BookingStruct) CreateBookingKafka(ctx context.Context, booking api.BookingRequest) (interface{}, error) {
	sentUuid := uuid.New().String()
	kafkaBooking := api.KafkaBookingRequest{
		UUID:          sentUuid,
		CustomerID:    booking.CustomerID,
		BarberId:      booking.BarberId,
		ResultId:      booking.ResultId,
		Status:        booking.Status,
		Price:         booking.Price,
		SlotId:        booking.SlotId,
		FeedBackId:    booking.FeedBackId,
		ListServiceId: booking.ListServiceId,
	}
	serializedBookingRequest, err := json.Marshal(kafkaBooking)
	if err != nil {
		b.logger.Log("Failed to serialize booking request: %s\n", err)
		return nil, err
	}
	err = kafka2.ProduceMessage(b.kafka, "booking", serializedBookingRequest)
	if err != nil {
		b.logger.Log("Failed to produce message: %s\n", err)
		return nil, err
	}
	fmt.Println("Message produced successfully!")

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBroker,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		return nil, err
	}
	defer c.Close()
	// Subscribe to the Kafka topic
	err = c.SubscribeTopics([]string{replyTopic}, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topic: %s\n", err)
		return nil, err
	}
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	run := true
	for run == true {
		fmt.Println("=========================Waiting for response=======================")
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
				var booking api.KafkaBookingResponse
				err := json.Unmarshal(e.Value, &booking)
				if err != nil {
					fmt.Printf("Failed to deserialize message: %s\n", err)
					return booking, nil
				}
			case kafka.Error:
				// Handle Kafka errors
				fmt.Printf("Error: %v\n", e)
				return nil, e
			}
		}
	}

	return nil, nil
}

func (b BookingStruct) CreateBooking(ctx context.Context, booking api.BookingRequest) (interface{}, error) {
	//call grpc api
	client := pb.NewUserServiceClient(b.connAccount)
	clientTimeslot := pbTimeslot.NewTimeslotServiceClient(b.connTimeslot)
	clientService := pbService.NewServicingServiceClient(b.connService)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	checkBarber, err := client.CheckExistedBarber(ctx, &pb.CheckExistedBarberRequest{Id: int32(booking.BarberId)})
	if err != nil {
		fmt.Printf("error when checking account: %v", err)
		return nil, err
	}
	if checkBarber.Value != "BARBER" {
		fmt.Println("barber id is not valid")
		return nil, errors.New("barber id is not valid")
	}

	checkUser, err := client.CheckExistedBarber(ctx, &pb.CheckExistedBarberRequest{Id: int32(booking.CustomerID)})
	if err != nil {
		fmt.Printf("error when checking account: %v", err)
		return nil, err
	}

	fmt.Printf("checkBarber: %s \n ", checkBarber.Value)
	fmt.Printf("checkUser: %s \n", checkUser.Value)

	checkTimeslot, err := clientTimeslot.CheckAvailableTimeslot(ctx, &pbTimeslot.CheckAvailableTimeslotRequest{Id: int32(booking.SlotId)})
	if err != nil {
		fmt.Printf("error when checking timeslot: %v", err)
		return nil, err
	}
	if checkTimeslot.Status != "Available" {
		fmt.Println("timeslot is not available")
		return nil, errors.New("timeslot is not available")
	}

	for _, service := range booking.ListServiceId {
		checkService, err := clientService.GetServiceById(ctx, &pbService.GetServiceByIdRequest{Id: int32(service)})
		if err != nil || checkService == nil {
			fmt.Printf("error when checking service: %v", err)
			return nil, errors.New("error when checking service")
		}
	}

	updatedTimeslot, err := clientTimeslot.UpdateStatusTimeslot(ctx, &pbTimeslot.UpdateStatusTimeslotRequest{Id: int32(booking.SlotId), Status: "Booked"})
	if err != nil {
		fmt.Printf("error when updating timeslot: %v", err)
		return nil, errors.New("error when updating timeslot")
	}

	fmt.Printf("updatedTimeslot successfully: %v \n", updatedTimeslot)

	resp, err := b.repository.CreateBooking(ctx, booking)
	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}
	err = b.repository.CreateBookingDetail(ctx, booking.ListServiceId, resp.ID)
	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}

	return resp, nil
}

func (b BookingStruct) GetBooking(ctx context.Context, id string) (interface{}, error) {
	clientService := pbService.NewServicingServiceClient(b.connService)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	idValue, _ := strconv.Atoi(id)
	resp, err := b.repository.GetBookingById(ctx, idValue)
	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}
	listIdServices, err := b.repository.GetListIdServiceByBookingId(ctx, idValue)
	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}
	var listServiceName []string
	for _, id := range listIdServices {
		service, err := clientService.GetServiceById(ctx, &pbService.GetServiceByIdRequest{Id: int32(id)})
		if err != nil {
			fmt.Printf("error when getting service: %v", err)
			return nil, err
		}
		listServiceName = append(listServiceName, service.Name)
	}

	return api.BookingResponse{
		ID:           resp.ID,
		CustomerID:   resp.CustomerID,
		BarberId:     resp.BarberId,
		ResultId:     resp.ResultId,
		Status:       resp.Status,
		Price:        resp.Price,
		SlotId:       resp.SlotId,
		FeedBackId:   resp.FeedBackId,
		CreatedAt:    resp.CreatedAt,
		UpdatedAt:    resp.UpdatedAt,
		ListServices: listServiceName,
	}, nil
}

func (b BookingStruct) GetListBooking(ctx context.Context) (interface{}, error) {
	resp, err := b.repository.GetListBooking(ctx)
	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}
	return resp, nil
}

func (b BookingStruct) UpdateBooking(ctx context.Context, booking api.UpdateBookingRequest) (interface{}, error) {
	_, err := b.repository.GetBookingById(ctx, booking.Id)
	if err != nil {
		fmt.Printf("booking not found")
		return nil, err
	}

	updatedBooking, err := b.repository.UpdateBooking(ctx, booking)
	if err != nil {
		fmt.Printf("booking not found")
		return nil, err
	}
	return updatedBooking, nil
}
