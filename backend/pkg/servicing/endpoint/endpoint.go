package endpoint

import (
	"DoAn/pkg/servicing"
	"DoAn/pkg/servicing/api"
	"context"
	"github.com/go-kit/kit/endpoint"
	"strconv"
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
