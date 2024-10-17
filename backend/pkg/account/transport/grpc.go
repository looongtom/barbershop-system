package transport

import (
	"DoAn/endpoint"
	"DoAn/entity"
	"DoAn/pb"
	"context"

	"google.golang.org/protobuf/types/known/wrapperspb"

	gt "github.com/go-kit/kit/transport/grpc"
)

type GRPCServer struct {
	checkExistedBarber gt.Handler
	verifyToken        gt.Handler
	getAccById         gt.Handler
}

func (G GRPCServer) GetAccountById(ctx context.Context, request *pb.CheckExistedBarberRequest) (*pb.Account, error) {
	_, resp, err := G.getAccById.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	var resGrpc *entity.Account
	resGrpc = resp.(*entity.Account)
	return &pb.Account{
		Id:          int32(resGrpc.ID),
		Username:    resGrpc.Username,
		Email:       resGrpc.Email,
		Role:        int32(resGrpc.Role),
		Phonenumber: resGrpc.PhoneNumber,
		Fullname:    resGrpc.FullName,
		About:       resGrpc.About,
		Avatar:      resGrpc.Avatar,
		CreatedAt:   resGrpc.CreatedAt,
		UpdatedAt:   resGrpc.UpdatedAt,
	}, nil
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
		getAccById: gt.NewServer(
			endpoints.GetAccount,
			endpoint.DecodeCheckExistedUser,
			endpoint.EncodeGetAccount,
		),
	}
}
