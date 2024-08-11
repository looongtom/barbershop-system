package transport

import (
	"DoAn/pkg/servicing/endpoint"
	"DoAn/pkg/servicing/pb"
	"context"

	gt "github.com/go-kit/kit/transport/grpc"
)

type GRPCServer struct {
	getServicingByBookingId gt.Handler
}

func (G GRPCServer) GetServiceById(ctx context.Context, request *pb.GetServiceByIdRequest) (*pb.Servicing, error) {
	_, resp, err := G.getServicingByBookingId.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.Servicing), nil
}

func NewGRPCServer(_ context.Context, endpoints endpoint.Endpoints) pb.ServicingServiceServer {
	return &GRPCServer{
		getServicingByBookingId: gt.NewServer(
			endpoints.GetServiceById,
			endpoint.DecodeGetServiceById,
			endpoint.EncodeGetServiceById,
		),
	}
}
