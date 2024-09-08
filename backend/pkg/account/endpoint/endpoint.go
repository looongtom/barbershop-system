package endpoint

import (
	"DoAn/pkg/account"
	"DoAn/pkg/account/api"
	"DoAn/pkg/account/auth"
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CheckExistedBarberEndpoint endpoint.Endpoint
	VerifyTokenEndpoint        endpoint.Endpoint
}

func MakeCheckExistedBarberEndpoint(s account.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(api.CheckExistedBarberRequest)
		return s.CheckExistedBarber(ctx, req.UserId)
	}
}

func MakeVerifyTokenEndpoint(svc auth.AuthenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(api.VerifyTokenRequest)
		return svc.VerifyToken(ctx, req.Token)
	}
}
