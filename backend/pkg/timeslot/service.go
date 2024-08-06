package timeslot

import (
	"DoAn/pkg/timeslot/api"
	"context"
)

type TimeSlotService interface {
	CreateTimeSlot(ctx context.Context, timeslot api.CreateOrUpdateTimeslotRequest) (interface{}, error)
	GetListTimeSlotByBarberId(ctx context.Context, findTimeSlot api.FindTimeslotRequest) (interface{}, error)

	CreateListTimeSlot(ctx context.Context, timeslots []api.CreateOrUpdateTimeslotRequest) (interface{}, error)
	//UpdateTimeSlot(ctx context.Context, timeslot api.CreateOrUpdateTimeslotRequest) (interface{}, error)
}
