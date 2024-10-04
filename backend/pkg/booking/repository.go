package booking

import (
	"DoAn/api"
	"DoAn/entity"
	"context"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking api.BookingRequest) (entity.Booking, error)
	GetListBooking(ctx context.Context) ([]entity.Booking, error)
	GetBookingById(ctx context.Context, id int) (entity.Booking, error)
	UpdateBooking(ctx context.Context, booking api.UpdateBookingRequest) (entity.Booking, error)
	FindBookingByUserOrBarber(ctx context.Context, findReq api.FindBookingRequest) ([]entity.Booking, error)

	CreateBookingDetail(ctx context.Context, listService []int, bookingId int) error
	GetListIdServiceByBookingId(ctx context.Context, id int) ([]int, error)
}
