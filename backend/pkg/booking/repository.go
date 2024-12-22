package booking

import (
	"context"

	"DoAn/api"
	"DoAn/entity"
	"DoAn/mapper"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking api.BookingRequest) (entity.Booking, error)
	GetListBooking(ctx context.Context, page, pageSize int) ([]mapper.BookingMapper, error)
	GetTotalCountBooking(ctx context.Context) (int, error)
	GetTotalCountBookingByUserId(ctx context.Context, id int) (int, error)
	GetTotalCountBookingByUserIdAndBookedDate(ctx context.Context, id int, bookedDate string) (int, error)
	GetTotalCountBookingByBarberId(ctx context.Context, id int) (int, error)
	GetTotalCountBookingByBarberIdAndBookedDate(ctx context.Context, id int, bookedDate string) (int, error)
	GetBookingById(ctx context.Context, id int) (entity.Booking, error)

	UpdateBooking(ctx context.Context, booking api.UpdateBookingRequest) (entity.Booking, error)
	UpdateBookingDetailService(ctx context.Context, listService []int, bookingId int) error

	FindBookingByUser(ctx context.Context, findReq api.FindListBookingRequest) ([]mapper.BookingMapper, error)
	FindBookingByUserAndBookedDate(ctx context.Context, findReq api.FindListBookingRequest) ([]mapper.BookingMapper, error)
	FindBookingByBarberAndBookedDate(ctx context.Context, findReq api.FindListBookingRequest) ([]mapper.BookingMapper, error)
	FindBookingByBarber(ctx context.Context, findReq api.FindListBookingRequest) ([]mapper.BookingMapper, error)

	CreateBookingDetail(ctx context.Context, listService []int, bookingId int) error
	GetListIdServiceByBookingId(ctx context.Context, id int) ([]int, error)
}
