package service

import (
	"DoAn/mapper"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"DoAn"
	"DoAn/api"
	kafka2 "DoAn/kafka"
	"DoAn/pb"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
)

const (
	bookingTopic = "booking"
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

func (b BookingStruct) FindBookingByBarber(ctx context.Context, findReq api.FindListBookingRequest) (int, interface{}, error) {
	resp, err := b.repository.FindBookingByBarber(ctx, findReq)
	if err != nil {
		return 0, nil, err
	}
	totalCount, err := b.repository.GetTotalCountBookingByBarberId(ctx, findReq.Account)
	if err != nil {
		return 0, nil, err
	}
	clientService := pb.NewServicingServiceClient(b.connService)
	clientTimeslotSvc := pb.NewTimeslotServiceClient(b.connTimeslot)
	clientAccountSvc := pb.NewUserServiceClient(b.connAccount)
	for i, booking := range resp {
		for _, id := range booking.ListServices {
			service, err := clientService.GetServiceById(ctx, &pb.GetServiceByIdRequest{Id: int32(id)})
			if err != nil {
				fmt.Printf("error when getting service: %v", err)
				return 0, nil, err
			}
			resp[i].ListServiceStruct = append(resp[i].ListServiceStruct, mapper.BookingService{
				ID:          int(service.Id),
				Name:        service.Name,
				Price:       int(service.Price),
				Description: service.Description,
				Url:         service.Url,
			})
		}
		timeslotInfo, err := clientTimeslotSvc.CheckAvailableTimeslot(ctx, &pb.CheckAvailableTimeslotRequest{Id: int32(booking.SlotId)})
		if err != nil {
			fmt.Printf("error when getting timeslot: %v", err)
			return 0, nil, err
		}
		resp[i].TimeSlot = mapper.BookingTimeSlot{
			ID:         int(timeslotInfo.Id),
			StartTime:  timeslotInfo.StartTime,
			BookedDate: timeslotInfo.BookedDate,
		}
		barberInfo, err := clientAccountSvc.GetAccountById(ctx, &pb.CheckExistedBarberRequest{Id: int32(booking.BarberId)})
		if err != nil {
			fmt.Printf("error when getting barber: %v", err)
			return 0, nil, err
		}
		customerInfo, err := clientAccountSvc.GetAccountById(ctx, &pb.CheckExistedBarberRequest{Id: int32(booking.CustomerID)})
		if err != nil {
			fmt.Printf("error when getting customer: %v", err)
			return 0, nil, err
		}
		resp[i].BarberName = barberInfo.Fullname
		resp[i].CustomerName = customerInfo.Fullname
	}
	return totalCount, resp, nil
}

func (b BookingStruct) FindBookingByUser(ctx context.Context, findReq api.FindListBookingRequest) (int, interface{}, error) {
	resp, err := b.repository.FindBookingByUser(ctx, findReq)
	if err != nil {
		return 0, nil, err
	}
	totalCount, err := b.repository.GetTotalCountBookingByUserId(ctx, findReq.Account)
	if err != nil {
		return 0, nil, err
	}
	clientService := pb.NewServicingServiceClient(b.connService)
	clientTimeslotSvc := pb.NewTimeslotServiceClient(b.connTimeslot)
	clientAccountSvc := pb.NewUserServiceClient(b.connAccount)
	for i, booking := range resp {
		for _, id := range booking.ListServices {
			service, err := clientService.GetServiceById(ctx, &pb.GetServiceByIdRequest{Id: int32(id)})
			if err != nil {
				fmt.Printf("error when getting service: %v", err)
				return 0, nil, err
			}
			resp[i].ListServiceStruct = append(resp[i].ListServiceStruct, mapper.BookingService{
				ID:          int(service.Id),
				Name:        service.Name,
				Price:       int(service.Price),
				Description: service.Description,
				Url:         service.Url,
			})
		}
		timeslotInfo, err := clientTimeslotSvc.CheckAvailableTimeslot(ctx, &pb.CheckAvailableTimeslotRequest{Id: int32(booking.SlotId)})
		if err != nil {
			fmt.Printf("error when getting timeslot: %v", err)
			return 0, nil, err
		}
		resp[i].TimeSlot = mapper.BookingTimeSlot{
			ID:         int(timeslotInfo.Id),
			StartTime:  timeslotInfo.StartTime,
			BookedDate: timeslotInfo.BookedDate,
		}
		barberInfo, err := clientAccountSvc.GetAccountById(ctx, &pb.CheckExistedBarberRequest{Id: int32(booking.BarberId)})
		if err != nil {
			fmt.Printf("error when getting barber: %v", err)
			return 0, nil, err
		}
		customerInfo, err := clientAccountSvc.GetAccountById(ctx, &pb.CheckExistedBarberRequest{Id: int32(booking.CustomerID)})
		if err != nil {
			fmt.Printf("error when getting customer: %v", err)
			return 0, nil, err
		}
		resp[i].BarberName = barberInfo.Fullname
		resp[i].CustomerName = customerInfo.Fullname
	}
	return totalCount, resp, nil
}

func (b BookingStruct) CreateBookingKafka(ctx context.Context, booking api.BookingRequest) (interface{}, error) {
	clientService := pb.NewServicingServiceClient(b.connService)
	clientAccountSvc := pb.NewUserServiceClient(b.connAccount)
	clientTimeslotSvc := pb.NewTimeslotServiceClient(b.connTimeslot)
	kafkaBooking := api.BookingRequest{
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
	err = kafka2.ProduceMessage(b.kafka, bookingTopic, serializedBookingRequest)
	if err != nil {
		b.logger.Log("Failed to produce message: %s\n", err)
		return nil, err
	}
	barberInfo, err := clientAccountSvc.GetAccountById(ctx, &pb.CheckExistedBarberRequest{Id: int32(booking.BarberId)})
	if err != nil {
		fmt.Printf("error when getting barber: %v", err)
		return nil, err
	}
	customerInfo, err := clientAccountSvc.GetAccountById(ctx, &pb.CheckExistedBarberRequest{Id: int32(booking.CustomerID)})
	if err != nil {
		fmt.Printf("error when getting customer: %v", err)
		return nil, err
	}
	timeslotInfo, err := clientTimeslotSvc.CheckAvailableTimeslot(ctx, &pb.CheckAvailableTimeslotRequest{Id: int32(booking.SlotId)})
	if err != nil {
		fmt.Printf("error when getting timeslot: %v", err)
		return nil, err
	}
	var listServiceName []api.ServiceResponse
	for _, id := range booking.ListServiceId {
		service, err := clientService.GetServiceById(ctx, &pb.GetServiceByIdRequest{Id: int32(id)})
		if err != nil {
			fmt.Printf("error when getting service: %v", err)
			return nil, err
		}
		listServiceName = append(listServiceName, api.ServiceResponse{
			ID:          int(service.Id),
			Name:        service.Name,
			Description: service.Description,
			Price:       int(service.Price),
			Url:         service.Url,
		})
	}

	return api.BookingResponse{
		CustomerID:   int(customerInfo.Id),
		CustomerName: customerInfo.Fullname,
		BarberId:     int(barberInfo.Id),
		BarberName:   barberInfo.Fullname,
		ResultId:     booking.ResultId,
		Status:       booking.Status,
		Price:        booking.Price,
		SlotId:       booking.SlotId,
		BookedDate:   timeslotInfo.BookedDate,
		StartTime:    timeslotInfo.StartTime,
		FeedBackId:   booking.FeedBackId,
		ListServices: listServiceName,
	}, nil
}

func (b BookingStruct) CreateBooking(ctx context.Context, booking api.BookingRequest) (interface{}, error) {
	// call grpc api
	client := pb.NewUserServiceClient(b.connAccount)
	clientTimeslot := pb.NewTimeslotServiceClient(b.connTimeslot)
	clientService := pb.NewServicingServiceClient(b.connService)
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

	checkTimeslot, err := clientTimeslot.CheckAvailableTimeslot(ctx, &pb.CheckAvailableTimeslotRequest{Id: int32(booking.SlotId)})
	if err != nil {
		fmt.Printf("error when checking timeslot: %v", err)
		return nil, err
	}
	if checkTimeslot.Status != "Available" {
		fmt.Println("timeslot is not available")
		return nil, errors.New("timeslot is not available")
	}

	for _, service := range booking.ListServiceId {
		checkService, err := clientService.GetServiceById(ctx, &pb.GetServiceByIdRequest{Id: int32(service)})
		if err != nil || checkService == nil {
			fmt.Printf("error when checking service: %v", err)
			return nil, errors.New("error when checking service")
		}
	}

	updatedTimeslot, err := clientTimeslot.UpdateStatusTimeslot(ctx, &pb.UpdateStatusTimeslotRequest{Id: int32(booking.SlotId), Status: "Booked"})
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
	//serializedBooking, err := json.Marshal(resp)
	//if err != nil {
	//	b.logger.Log("Failed to serialize booking request: %s\n", err)
	//	return nil, err
	//}
	//err = kafka2.ProduceMessage(b.kafka, bookingTopic, serializedBooking)
	//if err != nil {
	//	b.logger.Log("Failed to produce message: %s\n", err)
	//	return nil, err
	//}

	err = b.repository.CreateBookingDetail(ctx, booking.ListServiceId, resp.ID)
	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}

	return resp, nil
}

func (b BookingStruct) GetBooking(ctx context.Context, id string) (interface{}, error) {
	clientService := pb.NewServicingServiceClient(b.connService)
	clientAccountSvc := pb.NewUserServiceClient(b.connAccount)
	clientTimeslotSvc := pb.NewTimeslotServiceClient(b.connTimeslot)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	idValue, _ := strconv.Atoi(id)
	resp, err := b.repository.GetBookingById(ctx, idValue)
	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}
	barberInfo, err := clientAccountSvc.GetAccountById(ctx, &pb.CheckExistedBarberRequest{Id: int32(resp.BarberId)})
	if err != nil {
		fmt.Printf("error when getting barber: %v", err)
		return nil, err
	}
	customerInfo, err := clientAccountSvc.GetAccountById(ctx, &pb.CheckExistedBarberRequest{Id: int32(resp.CustomerID)})
	if err != nil {
		fmt.Printf("error when getting customer: %v", err)
		return nil, err
	}
	timeslotInfo, err := clientTimeslotSvc.CheckAvailableTimeslot(ctx, &pb.CheckAvailableTimeslotRequest{Id: int32(resp.SlotId)})
	if err != nil {
		fmt.Printf("error when getting timeslot: %v", err)
		return nil, err
	}
	listIdServices, err := b.repository.GetListIdServiceByBookingId(ctx, idValue)
	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}
	var listServiceName []api.ServiceResponse
	for _, id := range listIdServices {
		service, err := clientService.GetServiceById(ctx, &pb.GetServiceByIdRequest{Id: int32(id)})
		if err != nil {
			fmt.Printf("error when getting service: %v", err)
			return nil, err
		}
		listServiceName = append(listServiceName, api.ServiceResponse{
			ID:          int(service.Id),
			Name:        service.Name,
			Description: service.Description,
			Price:       int(service.Price),
			Url:         service.Url,
		})
	}

	return api.BookingResponse{
		ID:           resp.ID,
		CustomerID:   resp.CustomerID,
		CustomerName: customerInfo.Fullname,
		BarberId:     resp.BarberId,
		BarberName:   barberInfo.Fullname,
		ResultId:     resp.ResultId,
		Status:       resp.Status,
		Price:        resp.Price,
		SlotId:       resp.SlotId,
		BookedDate:   timeslotInfo.BookedDate,
		StartTime:    timeslotInfo.StartTime,
		FeedBackId:   resp.FeedBackId,
		CreatedAt:    resp.CreatedAt,
		UpdatedAt:    resp.UpdatedAt,
		ListServices: listServiceName,
	}, nil
}

func (b BookingStruct) GetListBooking(ctx context.Context, page, pageSize int) (int, interface{}, error) {
	resp, err := b.repository.GetListBooking(ctx, page, pageSize)
	if err != nil {
		errMsg := err.Error()
		return 0, errMsg, err
	}
	clientService := pb.NewServicingServiceClient(b.connService)
	clientTimeslotSvc := pb.NewTimeslotServiceClient(b.connTimeslot)
	clientAccountSvc := pb.NewUserServiceClient(b.connAccount)

	for i, booking := range resp {
		for _, id := range booking.ListServices {
			service, err := clientService.GetServiceById(ctx, &pb.GetServiceByIdRequest{Id: int32(id)})
			if err != nil {
				fmt.Printf("error when getting service: %v", err)
				return 0, nil, err
			}
			resp[i].ListServiceStruct = append(resp[i].ListServiceStruct, mapper.BookingService{
				ID:          int(service.Id),
				Name:        service.Name,
				Price:       int(service.Price),
				Description: service.Description,
				Url:         service.Url,
			})
		}
		timeslotInfo, err := clientTimeslotSvc.CheckAvailableTimeslot(ctx, &pb.CheckAvailableTimeslotRequest{Id: int32(booking.SlotId)})
		if err != nil {
			fmt.Printf("error when getting timeslot: %v", err)
			return 0, nil, err
		}
		resp[i].TimeSlot = mapper.BookingTimeSlot{
			ID:         int(timeslotInfo.Id),
			StartTime:  timeslotInfo.StartTime,
			BookedDate: timeslotInfo.BookedDate,
		}
		barberInfo, err := clientAccountSvc.GetAccountById(ctx, &pb.CheckExistedBarberRequest{Id: int32(booking.BarberId)})
		if err != nil {
			fmt.Printf("error when getting barber: %v", err)
			return 0, nil, err
		}
		customerInfo, err := clientAccountSvc.GetAccountById(ctx, &pb.CheckExistedBarberRequest{Id: int32(booking.CustomerID)})
		if err != nil {
			fmt.Printf("error when getting customer: %v", err)
			return 0, nil, err
		}
		resp[i].BarberName = barberInfo.Fullname
		resp[i].CustomerName = customerInfo.Fullname
	}
	totalCount, err := b.repository.GetTotalCountBooking(ctx)
	if err != nil {
		return 0, nil, err
	}
	return totalCount, resp, nil
}

func compareListService(listIdServices []int, listServiceId []int) bool {
	if len(listIdServices) != len(listServiceId) {
		return false
	}
	for _, id := range listServiceId {
		if !contains(listIdServices, id) {
			return false
		}
	}
	return true
}

func contains(services []int, id int) bool {
	for _, service := range services {
		if service == id {
			return true
		}
	}
	return false
}

func (b BookingStruct) UpdateBookingService(ctx context.Context, booking api.UpdateBookingServiceRequest) error {
	clientService := pb.NewServicingServiceClient(b.connService)
	for _, id := range booking.ListServiceId {
		_, err := clientService.GetServiceById(ctx, &pb.GetServiceByIdRequest{Id: int32(id)})
		if err != nil {
			fmt.Printf("error when getting service: %v", err)
			return err
		}
	}

	err := b.repository.UpdateBookingDetailService(ctx, booking.ListServiceId, booking.Id)
	if err != nil {
		fmt.Printf("error when updating booking detail service: %v", err)
		return err
	}
	return nil
}

func (b BookingStruct) UpdateBooking(ctx context.Context, booking api.UpdateBookingRequest) (interface{}, error) {
	clientService := pb.NewServicingServiceClient(b.connService)
	clientAccountSvc := pb.NewUserServiceClient(b.connAccount)
	clientTimeslotSvc := pb.NewTimeslotServiceClient(b.connTimeslot)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := b.repository.GetBookingById(ctx, booking.Id)
	if err != nil {
		fmt.Printf("booking not found")
		return nil, err
	}

	if booking.SlotId != resp.SlotId {
		checkTimeslot, err := clientTimeslotSvc.CheckAvailableTimeslot(ctx, &pb.CheckAvailableTimeslotRequest{Id: int32(booking.SlotId)})
		if err != nil {
			fmt.Printf("error when checking timeslot: %v", err)
			return nil, err
		}
		if checkTimeslot.Status != "Available" {
			fmt.Println("timeslot is not available")
			return nil, errors.New("timeslot is not available")
		}
		updatedTimeslot, err := clientTimeslotSvc.UpdateStatusTimeslot(ctx, &pb.UpdateStatusTimeslotRequest{Id: int32(booking.SlotId), Status: "Booked"})
		if err != nil {
			fmt.Printf("error when updating timeslot: %v", err)
			return nil, err
		}
		b.logger.Log("updatedTimeslot successfully: %v \n", updatedTimeslot)

		updatedTimeslot, err = clientTimeslotSvc.UpdateStatusTimeslot(ctx, &pb.UpdateStatusTimeslotRequest{Id: int32(resp.SlotId), Status: "Available"})
		if err != nil {
			fmt.Printf("error when updating timeslot: %v", err)
			return nil, err
		}
		b.logger.Log(fmt.Sprintf("updatedTimeslot successfully: %v \n", updatedTimeslot))

	}

	updatedBooking, err := b.repository.UpdateBooking(ctx, booking)
	if err != nil {
		fmt.Printf("booking not found")
		return nil, err
	}

	listIdServices, err := b.repository.GetListIdServiceByBookingId(ctx, resp.ID)
	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}

	//compare listIdServices and booking.ListServiceId
	//if not equal, update booking detail service
	if len(listIdServices) != len(booking.ListServiceId) {
		err = b.repository.UpdateBookingDetailService(ctx, booking.ListServiceId, booking.Id)
		if err != nil {
			fmt.Printf("error when updating booking detail service: %v", err)
			return nil, err
		}
	} else {
		//check if listIdServices and booking.ListServiceId are equal
		//if not equal, update booking detail service
		if !compareListService(listIdServices, booking.ListServiceId) {
			err = b.repository.UpdateBookingDetailService(ctx, booking.ListServiceId, booking.Id)
			if err != nil {
				fmt.Printf("error when updating booking detail service: %v", err)
				return nil, err
			}
		}
	}
	barberInfo, err := clientAccountSvc.GetAccountById(ctx, &pb.CheckExistedBarberRequest{Id: int32(resp.BarberId)})
	if err != nil {
		fmt.Printf("error when getting barber: %v", err)
		return nil, err
	}
	customerInfo, err := clientAccountSvc.GetAccountById(ctx, &pb.CheckExistedBarberRequest{Id: int32(resp.CustomerID)})
	if err != nil {
		fmt.Printf("error when getting customer: %v", err)
		return nil, err
	}
	timeslotInfo, err := clientTimeslotSvc.CheckAvailableTimeslot(ctx, &pb.CheckAvailableTimeslotRequest{Id: int32(resp.SlotId)})
	if err != nil {
		fmt.Printf("error when getting timeslot: %v", err)
		return nil, err
	}

	listIdServicesDisplay, err := b.repository.GetListIdServiceByBookingId(ctx, updatedBooking.ID)
	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}

	var listServiceName []api.ServiceResponse
	for _, id := range listIdServicesDisplay {
		service, err := clientService.GetServiceById(ctx, &pb.GetServiceByIdRequest{Id: int32(id)})
		if err != nil {
			fmt.Printf("error when getting service: %v", err)
			return nil, err
		}
		listServiceName = append(listServiceName, api.ServiceResponse{
			ID:          int(service.Id),
			Name:        service.Name,
			Description: service.Description,
			Price:       int(service.Price),
			Url:         service.Url,
		})
	}

	return api.BookingResponse{
		ID:           updatedBooking.ID,
		CustomerID:   updatedBooking.CustomerID,
		CustomerName: customerInfo.Fullname,
		BarberId:     updatedBooking.BarberId,
		BarberName:   barberInfo.Fullname,
		ResultId:     updatedBooking.ResultId,
		Status:       updatedBooking.Status,
		Price:        updatedBooking.Price,
		SlotId:       updatedBooking.SlotId,
		BookedDate:   timeslotInfo.BookedDate,
		StartTime:    timeslotInfo.StartTime,
		FeedBackId:   updatedBooking.FeedBackId,
		CreatedAt:    updatedBooking.CreatedAt,
		UpdatedAt:    updatedBooking.UpdatedAt,
		ListServices: listServiceName,
	}, nil
}
