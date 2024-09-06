package endpoint

import (
	"DoAn/pkg/servicing/api"
	"DoAn/pkg/servicing/entity"
	"DoAn/pkg/servicing/pb"
	"context"
)

func DecodeGetServiceById(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetServiceByIdRequest)
	return api.GetServiceByIdRequest{
		Id: int(req.Id),
	}, nil
}

func EncodeGetServiceById(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*entity.Servicing)
	return &pb.Servicing{
		Id:          int32(res.ID),
		Name:        res.Name,
		Price:       int32(res.Price),
		Description: res.Description,
		Url:         res.Url,
		CategoryId:  int32(res.CategoryID),
	}, nil
}
