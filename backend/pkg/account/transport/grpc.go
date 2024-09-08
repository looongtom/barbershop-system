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
	verifyToken        gt.Handler
}

func (G GRPCServer) VerifyToken(ctx context.Context, request *pb.VerifyTokenRequest) (*wrapperspb.StringValue, error) {
	_, resp, err := G.verifyToken.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	var resGrpc wrapperspb.StringValue
	resGrpc.Value = resp.(string)
	return &resGrpc, nil
}

func (G GRPCServer) CheckExistedBarber(ctx context.Context, request *pb.CheckExistedBarberRequest) (*wrapperspb.StringValue, error) {
	_, resp, err := G.checkExistedBarber.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	var resGrpc wrapperspb.StringValue
	resGrpc.Value = resp.(string)
	return &resGrpc, nil
}

func NewGRPCServer(_ context.Context, endpoints endpoint.Endpoints) pb.UserServiceServer {
	return &GRPCServer{
		checkExistedBarber: gt.NewServer(
			endpoints.CheckExistedBarberEndpoint,
			endpoint.DecodeCheckExistedUser,
			endpoint.EncodeCheckExistedUser,
		),
		verifyToken: gt.NewServer(
			endpoints.VerifyTokenEndpoint,
			endpoint.DecodeVerifyToken,
			endpoint.EncodeVerifyToken,
		),
	}
}
