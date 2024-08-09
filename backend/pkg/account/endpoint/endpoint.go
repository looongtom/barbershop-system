package endpoint

import (
	"DoAn/pkg/account"
	"context"
	"github.com/go-kit/kit/endpoint"
)

type CheckExistedBarberRequest struct {
	UserId int `json:"user_id"`
}

type Endpoints struct {
	CheckExistedBarberEndpoint endpoint.Endpoint
}

func MakeCheckExistedBarberEndpoint(s account.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CheckExistedBarberRequest)
		return s.CheckExistedBarber(ctx, req.UserId)
	}
}
