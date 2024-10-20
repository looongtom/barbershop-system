package booking

import (
	"DoAn/api"
	"context"
)

type BookingService interface {
	CreateBooking(ctx context.Context, booking api.BookingRequest) (interface{}, error)
	CreateBookingKafka(ctx context.Context, booking api.BookingRequest) (interface{}, error)
	GetBooking(ctx context.Context, id string) (interface{}, error)
	GetListBooking(ctx context.Context, page, pageSize int) (int, interface{}, error)
	UpdateBooking(ctx context.Context, booking api.UpdateBookingRequest) (interface{}, error)
	UpdateBookingService(ctx context.Context, booking api.UpdateBookingServiceRequest) error

	FindBookingByUser(ctx context.Context, findReq api.FindListBookingRequest) (int, interface{}, error)
	FindBookingByBarber(ctx context.Context, findReq api.FindListBookingRequest) (int, interface{}, error)
}
