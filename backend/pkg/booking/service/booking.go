package service

import (
	"DoAn/pkg/booking"
	"DoAn/pkg/booking/api"
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"strconv"
)

type BookingStruct struct {
	repository booking.BookingRepository
	logger     log.Logger
}

func NewService(repo booking.BookingRepository, logger log.Logger) booking.BookingService {
	return &BookingStruct{
		repository: repo,
		logger:     logger,
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
