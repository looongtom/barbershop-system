package timeslot

import (
	"DoAn/pkg/timeslot/api"
	"context"
)

type TimeSlotService interface {
	CreateOrUpdateTimeSlot(ctx context.Context, timeslot api.CreateOrUpdateTimeslotRequest) (interface{}, error)
	GetListTimeSlotByBarberId(ctx context.Context, findTimeSlot api.FindTimeslotRequest) (interface{}, error)

	CreateListTimeSlot(ctx context.Context, timeslots []api.CreateOrUpdateTimeslotRequest) (interface{}, error)
	CheckExistTimeslot(ctx context.Context, id int) (interface{}, error)
	UpdateStatusTimeSlot(ctx context.Context, id int, status string) (interface{}, error)
}
