package endpoint

import (
	"DoAn/entity"
	"DoAn/pkg/timeslot/api"
	"DoAn/pkg/timeslot/pb"
	"context"
)

func DecodeFindExistTimeslot(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.FindTimeslotRequest)
	return api.FindTimeslotRequest{
		BarberId:   int(req.BarberId),
		StartTime:  req.StartTime,
		BookedDate: req.BookedDate,
		Status:     req.Status,
	}, nil
}

func EncodeFindExistTimeslot(_ context.Context, r interface{}) (interface{}, error) {
	res := r.([]entity.Timeslot)
	return res, nil
}

func DecodeCheckExistTimeslot(ctx context.Context, i interface{}) (request interface{}, err error) {
	req := i.(*pb.CheckAvailableTimeslotRequest)
	return api.CheckExistTimeslotRequest{
		Id: int(req.Id),
	}, nil
}

func DecodeUpdateTimeslotStatus(ctx context.Context, i interface{}) (request interface{}, err error) {
	req := i.(*pb.UpdateStatusTimeslotRequest)
	return api.UpdateTimeslotRequest{
		ID:     int(req.Id),
		Status: req.Status,
	}, nil
}

func EncodeCheckExistTimeslot(ctx context.Context, i interface{}) (response interface{}, err error) {
	resp := i.(*entity.Timeslot)
	return &pb.Timeslot{
		Id:         int32(resp.ID),
		StartTime:  resp.StartTime,
		BookedDate: resp.BookedDate,
		Status:     resp.Status,
		BarberId:   int32(resp.BarberId),
		CreatedAt:  resp.CreatedAt,
		UpdatedAt:  resp.UpdatedAt,
	}, nil
}

func EncodeUpdateTimeslotStatus(ctx context.Context, i interface{}) (response interface{}, err error) {
	resp := i.(entity.Timeslot)
	return &pb.Timeslot{
		Id:         int32(resp.ID),
		StartTime:  resp.StartTime,
		BookedDate: resp.BookedDate,
		Status:     resp.Status,
		BarberId:   int32(resp.BarberId),
		CreatedAt:  resp.CreatedAt,
		UpdatedAt:  resp.UpdatedAt,
	}, nil
}
