package endpoint

import (
	"DoAn"
	"DoAn/api"
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateBooking endpoint.Endpoint
}

func MakeCreateBookingEndpoint(svc booking.BookingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(api.BookingRequest)
		return svc.CreateBooking(ctx, req)
	}
}
