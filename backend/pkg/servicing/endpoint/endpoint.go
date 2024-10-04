package endpoint

import (
	"DoAn"
	"DoAn/api"
	"context"
	"strconv"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetServiceById endpoint.Endpoint
}

func MakeGetServiceByIdEndpoint(svc servicing.ServicingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(api.GetServiceByIdRequest)
		return svc.GetServicing(ctx, strconv.Itoa(req.Id))
	}
}
