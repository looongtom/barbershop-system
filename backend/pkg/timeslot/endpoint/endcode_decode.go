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
	//var listTimeslot []entity.Timeslot
	//for _, v := range res.Timeslots {
	//	listTimeslot = append(listTimeslot, entity.Timeslot{
	//		ID:         int(v.Id),
	//		BarberId:   int(v.BarberId),
	//		StartTime:  v.StartTime,
	//		BookedDate: v.BookedDate,
	//		Status:     v.Status,
	//	})
	//}
	return res, nil
}
