package transport

import (
	"DoAn/endpoint"
	"DoAn/pb"

	"context"

	gt "github.com/go-kit/kit/transport/grpc"
)

type GRPCServer struct {
	createBooking gt.Handler
}

func (G GRPCServer) CreateBooking(ctx context.Context, request *pb.BookingRequest) (*pb.Booking, error) {
	_, resp, err := G.createBooking.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.Booking), nil
}

func NewGRPCServer(_ context.Context, endpoints endpoint.Endpoints) pb.BookingServiceServer {
	return &GRPCServer{
		createBooking: gt.NewServer(
			endpoints.CreateBooking,
			endpoint.DecodeCreateBooking,
			endpoint.EncodeCreateBooking,
		),
	}

}
