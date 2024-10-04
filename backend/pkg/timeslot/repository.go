package timeslot

import (
	"DoAn/api"
	"DoAn/entity"
	"context"
)

type TimeslotRepository interface {
	CreateTimeSlot(ctx context.Context, timeslot api.CreateOrUpdateTimeslotRequest) (entity.Timeslot, error)
	CreateListTimeSlot(ctx context.Context, timeslots []api.CreateOrUpdateTimeslotRequest) ([]entity.Timeslot, error)

	GetListTimeSlotByBarberId(ctx context.Context, findTimeSlot api.FindTimeslotRequest) ([]entity.Timeslot, error)
	UpdateTimeSlot(ctx context.Context, timeslot api.CreateOrUpdateTimeslotRequest) (entity.Timeslot, error)
	UpdateStatusTimeSlot(ctx context.Context, id int, status string) (entity.Timeslot, error)
	CheckAvailableTimeSlot(ctx context.Context, checkExist api.CheckExistedTimeslotRequest) (bool, error)
	CheckExistedTimeSlotById(ctx context.Context, id int) (*entity.Timeslot, error)
}
