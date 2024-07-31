package booking

import (
	"DoAn/pkg/booking/api"
	"context"
)

type BookingService interface {
	CreateBooking(ctx context.Context, booking api.BookingRequest) (interface{}, error)
	GetBooking(ctx context.Context, id string) (interface{}, error)
	GetListBooking(ctx context.Context) (interface{}, error)
	UpdateBooking(ctx context.Context, booking api.UpdateBookingRequest) (interface{}, error)

	FindBookingByUserOrBarber(ctx context.Context, findReq api.FindBookingRequest) (interface{}, error)
}
