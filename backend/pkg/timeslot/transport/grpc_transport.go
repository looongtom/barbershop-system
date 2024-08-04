package transport

import (
	"DoAn/entity"
	"DoAn/pkg/timeslot/endpoint"
	"DoAn/pkg/timeslot/pb"
	"context"

	gt "github.com/go-kit/kit/transport/grpc"
)

type GRPCServer struct {
	findExistedTimeslot gt.Handler
}

func (G GRPCServer) FindExistTimeslot(ctx context.Context, request *pb.FindTimeslotRequest) (*pb.TimeslotList, error) {
	_, resp, err := G.findExistedTimeslot.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	var listGrpcTimeslot []*pb.Timeslot
	for _, v := range resp.([]entity.Timeslot) {
		listGrpcTimeslot = append(listGrpcTimeslot, &pb.Timeslot{
			Id:         int32(v.ID),
			StartTime:  v.StartTime,
			BookedDate: v.BookedDate,
			Status:     v.Status,
			BarberId:   int32(v.BarberId),
			CreatedAt:  v.CreatedAt,
			UpdatedAt:  v.UpdatedAt,
		})
	}
	return &pb.TimeslotList{Timeslots: listGrpcTimeslot}, nil
}

func NewGRPCServer(_ context.Context, endpoints endpoint.Endpoints) pb.TimeslotServiceServer {
	return &GRPCServer{
		findExistedTimeslot: gt.NewServer(
			endpoints.FindExistTimeslotEndpoint,
			endpoint.DecodeFindExistTimeslot,
			endpoint.EncodeFindExistTimeslot,
		),
	}
}
