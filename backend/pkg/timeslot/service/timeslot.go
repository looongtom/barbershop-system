package service

import (
	"DoAn/entity"
	"DoAn/pkg/timeslot"
	"DoAn/pkg/timeslot/api"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
)

type TimeSlotService struct {
	repository timeslot.TimeslotRepository
	logger     log.Logger
}

func NewService(repo timeslot.TimeslotRepository, logger log.Logger) timeslot.TimeSlotService {
	return &TimeSlotService{
		repository: repo,
		logger:     logger,
	}
}

func (t TimeSlotService) CreateListTimeSlot(ctx context.Context, timeslots []api.CreateOrUpdateTimeslotRequest) (interface{}, error) {
	var listRes []entity.Timeslot
	for _, timeslot := range timeslots {
		ok, err := t.repository.CheckAvailableTimeSlot(ctx, api.CheckExistedTimeslotRequest{
			BarberId:   timeslot.BarberId,
			StartTime:  timeslot.StartTime,
			BookedDate: timeslot.BookedDate,
		})
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			t.logger.Log(fmt.Sprintf("error %v while checking available timeslot %v \n", err, timeslot))
			continue
		}
		if ok {
			t.logger.Log(fmt.Sprintf("error timeslot already exist: %v \n", err, timeslot))
			continue
		}
		res, err := t.repository.CreateTimeSlot(ctx, timeslot)
		if err != nil {
			return nil, err
		} else {
			listRes = append(listRes, res)
		}
	}
	if len(listRes) == 0 {
		return nil, errors.New("create timeslot failed")
	}
	return listRes, nil
}

func (t TimeSlotService) CreateTimeSlot(ctx context.Context, timeslot api.CreateOrUpdateTimeslotRequest) (interface{}, error) {
	if &timeslot.ID != nil && timeslot.ID != 0 {
		existedTimeslot, err := t.repository.CheckExistedTimeSlotById(ctx, timeslot.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if &existedTimeslot != nil {
			if timeslot.StartTime == "" || &timeslot.StartTime == nil {
				timeslot.StartTime = existedTimeslot.StartTime
			}
			if timeslot.BookedDate == "" || &timeslot.BookedDate == nil {
				timeslot.BookedDate = existedTimeslot.BookedDate
			}
			if timeslot.Status == "" || &timeslot.Status == nil {
				timeslot.Status = existedTimeslot.Status
			}
			if &timeslot.BarberId == nil {
				timeslot.BarberId = existedTimeslot.BarberId
			}
			ok, err := t.repository.CheckAvailableTimeSlot(ctx, api.CheckExistedTimeslotRequest{
				BarberId:   timeslot.BarberId,
				StartTime:  timeslot.StartTime,
				BookedDate: timeslot.BookedDate,
			})
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
			if ok {
				return nil, errors.New(api.ErrTimeslotAlreadyExist)
			}
			resp, err := t.repository.UpdateTimeSlot(ctx, timeslot)
			if err != nil {
				return nil, err
			}
			return resp, nil
		}
		return nil, errors.New("timeslot not found")
	} else {
		ok, err := t.repository.CheckAvailableTimeSlot(ctx, api.CheckExistedTimeslotRequest{
			BarberId:   timeslot.BarberId,
			StartTime:  timeslot.StartTime,
			BookedDate: timeslot.BookedDate,
		})
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if ok {
			return nil, errors.New(api.ErrTimeslotAlreadyExist)
		}
		res, err := t.repository.CreateTimeSlot(ctx, timeslot)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func (t TimeSlotService) GetListTimeSlotByBarberId(ctx context.Context, findTimeSlot api.FindTimeslotRequest) (interface{}, error) {
	res, err := t.repository.GetListTimeSlotByBarberId(ctx, findTimeSlot)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//func (t TimeSlotService) UpdateTimeSlot(ctx context.Context, timeslot api.CreateOrUpdateTimeslotRequest) (interface{}, error) {
//}
