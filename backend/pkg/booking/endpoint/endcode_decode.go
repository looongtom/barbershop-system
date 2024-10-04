package endpoint

import (
	"DoAn/api"
	"DoAn/entity"
	"DoAn/pb"
	"context"
)

func DecodeCreateBooking(ctx context.Context, i interface{}) (request interface{}, err error) {
	req := i.(*pb.BookingRequest)
	listServiceId := make([]int, len(req.ListServiceId))
	for i, v := range req.ListServiceId {
		listServiceId[i] = int(v)
	}
	return api.BookingRequest{
		CustomerID:    int(req.CustomerId),
		BarberId:      int(req.BarberId),
		Status:        req.Status,
		Price:         req.Price,
		SlotId:        int(req.SlotId),
		ListServiceId: listServiceId,
	}, nil
}

func EncodeCreateBooking(ctx context.Context, i interface{}) (response interface{}, err error) {
	resp := i.(entity.Booking)
	return &pb.Booking{
		Id:         int32(resp.ID),
		CustomerId: int32(resp.CustomerID),
		BarberId:   int32(resp.BarberId),
		Status:     resp.Status,
		Price:      resp.Price,
		SlotId:     int32(resp.SlotId),
		CreatedAt:  int32(resp.CreatedAt),
		UpdatedAt:  int32(resp.UpdatedAt),
	}, nil
}
