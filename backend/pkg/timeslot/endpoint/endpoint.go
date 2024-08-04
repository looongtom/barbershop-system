package endpoint

import (
	"DoAn/pkg/timeslot"
	"DoAn/pkg/timeslot/api"
	"context"
	"github.com/go-kit/kit/endpoint"
)

type FindExistTimeslotRequest struct {
	StartTime  string `json:"start_time"`
	BookedDate string `json:"booked_date"`
	Status     string `json:"status"`
	BarberId   int    `json:"barber_id"`
}

type Endpoints struct {
	FindExistTimeslotEndpoint endpoint.Endpoint
}

func MakeFindExistTimeslotEndpoint(svc timeslot.TimeSlotService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(api.FindTimeslotRequest)
		return svc.GetListTimeSlotByBarberId(ctx, req)
	}
}
