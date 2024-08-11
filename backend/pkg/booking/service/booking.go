package service

import (
	"DoAn/pkg/account/pb"
	"DoAn/pkg/booking"
	"DoAn/pkg/booking/api"
	pbTimeslot "DoAn/pkg/timeslot/pb"
	"errors"

	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
)

type BookingStruct struct {
	repository   booking.BookingRepository
	connAccount  *grpc.ClientConn
	connTimeslot *grpc.ClientConn
	logger       log.Logger
}

func NewService(repo booking.BookingRepository, logger log.Logger, conn *grpc.ClientConn, conn2 *grpc.ClientConn) booking.BookingService {
	return &BookingStruct{
		repository:   repo,
		logger:       logger,
		connAccount:  conn,
		connTimeslot: conn2,
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

func (b BookingStruct) CreateBooking(ctx context.Context, booking api.BookingRequest) (interface{}, error) {
	//call grpc api
	client := pb.NewUserServiceClient(b.connAccount)
	clientTimeslot := pbTimeslot.NewTimeslotServiceClient(b.connTimeslot)
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
	idValue, _ := strconv.Atoi(id)
	resp, err := b.repository.GetBookingById(ctx, idValue)

	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}
	return resp, nil
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
