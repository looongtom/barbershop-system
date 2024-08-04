package endpoint

import (
	"DoAn/pkg/account/pb"
	"context"
)

func DecodeCheckExistedUser(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.CheckExistedBarberRequest)
	return CheckExistedBarberRequest{BarberID: int(req.Id)}, nil
}

func EncodeCheckExistedUser(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(bool)
	return res, nil
}
