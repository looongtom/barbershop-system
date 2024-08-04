package transport

import (
	"DoAn/pkg/account/endpoint"
	"DoAn/pkg/account/pb"
	"context"
	"google.golang.org/protobuf/types/known/wrapperspb"

	gt "github.com/go-kit/kit/transport/grpc"
)

type GRPCServer struct {
	checkExistedBarber gt.Handler
}

func (G GRPCServer) CheckExistedBarber(ctx context.Context, request *pb.CheckExistedBarberRequest) (*wrapperspb.BoolValue, error) {
	_, resp, err := G.checkExistedBarber.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	var grpcResp *wrapperspb.BoolValue
	if resp.(bool) {
		grpcResp = wrapperspb.Bool(true)
	} else {
		grpcResp = wrapperspb.Bool(false)
	}

	return grpcResp, nil
}

func NewGRPCServer(_ context.Context, endpoints endpoint.Endpoints) pb.UserServiceServer {
	return &GRPCServer{
		checkExistedBarber: gt.NewServer(
			endpoints.CheckExistedBarberEndpoint,
			endpoint.DecodeCheckExistedUser,
			endpoint.EncodeCheckExistedUser,
		),
	}
}
