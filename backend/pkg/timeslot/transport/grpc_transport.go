package transport

import (
	"DoAn/pkg/timeslot/endpoint"
	"DoAn/pkg/timeslot/entity"
	"DoAn/pkg/timeslot/pb"
	"context"

	gt "github.com/go-kit/kit/transport/grpc"
)

type GRPCServer struct {
	findExistedTimeslot    gt.Handler
	checkAvailableTimeslot gt.Handler
	updateStatusTimeslot   gt.Handler
}

func (G GRPCServer) UpdateStatusTimeslot(ctx context.Context, request *pb.UpdateStatusTimeslotRequest) (*pb.Timeslot, error) {
	_, resp, err := G.updateStatusTimeslot.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.Timeslot), nil
}

func (G GRPCServer) CheckAvailableTimeslot(ctx context.Context, request *pb.CheckAvailableTimeslotRequest) (*pb.Timeslot, error) {
	_, resp, err := G.checkAvailableTimeslot.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.Timeslot), nil
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
		checkAvailableTimeslot: gt.NewServer(
			endpoints.CheckExistTimeslotEndpoint,
			endpoint.DecodeCheckExistTimeslot,
			endpoint.EncodeCheckExistTimeslot,
		),
		updateStatusTimeslot: gt.NewServer(
			endpoints.UpdateTimeslotStatusEndpoint,
			endpoint.DecodeUpdateTimeslotStatus,
			endpoint.EncodeUpdateTimeslotStatus,
		),
	}
}
