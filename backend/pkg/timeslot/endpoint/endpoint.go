package endpoint

import (
	"DoAn"
	"DoAn/api"
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	FindExistTimeslotEndpoint    endpoint.Endpoint
	CheckExistTimeslotEndpoint   endpoint.Endpoint
	UpdateTimeslotStatusEndpoint endpoint.Endpoint
}

func MakeFindExistTimeslotEndpoint(svc timeslot.TimeSlotService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(api.FindTimeslotRequest)
		return svc.GetListTimeSlotByBarberId(ctx, req)
	}
}

func MakeCheckExistTimeslotEndpoint(svc timeslot.TimeSlotService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(api.CheckExistTimeslotRequest)
		return svc.CheckExistTimeslot(ctx, req.Id)
	}
}

func MakeUpdateTimeslotStatusEndpoint(svc timeslot.TimeSlotService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(api.UpdateTimeslotRequest)
		return svc.UpdateStatusTimeSlot(ctx, req.ID, req.Status)
	}
}
