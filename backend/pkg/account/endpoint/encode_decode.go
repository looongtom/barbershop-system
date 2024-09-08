package endpoint

import (
	"DoAn/pkg/account/api"
	"DoAn/pkg/account/pb"
	"context"
)

func DecodeVerifyToken(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.VerifyTokenRequest)
	return api.VerifyTokenRequest{Token: req.Token}, nil
}

func EncodeVerifyToken(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(string)
	return res, nil
}

func DecodeCheckExistedUser(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.CheckExistedBarberRequest)
	return api.CheckExistedBarberRequest{UserId: int(req.Id)}, nil
}

func EncodeCheckExistedUser(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(string)
	return res, nil
}
