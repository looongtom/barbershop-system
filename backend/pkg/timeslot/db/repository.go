package db

import (
	"DoAn/entity"
	"DoAn/pkg/timeslot"
	"DoAn/pkg/timeslot/api"
	"context"
	"database/sql"
	"github.com/go-kit/kit/log"
	"time"
)

type repo struct {
	db     *sql.DB
	logger log.Logger
}

func NewRepository(db *sql.DB, logger log.Logger) (timeslot.TimeslotRepository, error) {
	return &repo{
		db:     db,
		logger: logger,
	}, nil
}

func (r repo) CreateTimeSlot(ctx context.Context, timeslot api.CreateOrUpdateTimeslotRequest) (entity.Timeslot, error) {
	newTimeSlot := entity.Timeslot{
		StartTime:  timeslot.StartTime,
		BookedDate: timeslot.BookedDate,
		Status:     timeslot.Status,
		BarberId:   timeslot.BarberId,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}
	query := `
		INSERT INTO timeslot(start_time, booked_date, status, barber_id, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id
	`
	err := r.db.QueryRow(query, newTimeSlot.StartTime, newTimeSlot.BookedDate, newTimeSlot.Status, newTimeSlot.BarberId, newTimeSlot.CreatedAt, newTimeSlot.UpdatedAt).Scan(&newTimeSlot.ID)
	if err != nil {
		r.logger.Log("error while inserting data")
		return entity.Timeslot{}, err
	}
	return newTimeSlot, nil
}

func (r repo) GetListTimeSlotByBarberId(ctx context.Context, findTimeSlot api.FindTimeslotRequest) ([]entity.Timeslot, error) {
	query := `
		SELECT id,start_time,booked_date,status,barber_id,created_at,updated_at FROM timeslot WHERE barber_id = $1
		or booked_date = $2
		or start_time = $3
		or status = $4
`
	params := []interface{}{findTimeSlot.BarberId, findTimeSlot.BookedDate, findTimeSlot.StartTime, findTimeSlot.Status}
	rows, err := r.db.Query(query, params...)

	if err != nil {
		r.logger.Log("error while querying")
		return nil, err
	}
	var listTimeSlot []entity.Timeslot
	for rows.Next() {
		var timeslot entity.Timeslot
		err = rows.Scan(&timeslot.ID, &timeslot.StartTime, &timeslot.BookedDate, &timeslot.Status, &timeslot.BarberId, &timeslot.CreatedAt, &timeslot.UpdatedAt)
		if err != nil {
			r.logger.Log("error while scanning")
			return nil, err
		}
		listTimeSlot = append(listTimeSlot, timeslot)
	}
	return listTimeSlot, nil
}

func (r repo) CheckExistedTimeSlotById(ctx context.Context, id int) (*entity.Timeslot, error) {
	query := `SELECT id,start_time,booked_date,status,barber_id,created_at,updated_at FROM timeslot WHERE id = $1`
	var timeslot entity.Timeslot
	err := r.db.QueryRow(query, id).Scan(&timeslot.ID, &timeslot.StartTime, &timeslot.BookedDate, &timeslot.Status, &timeslot.BarberId, &timeslot.CreatedAt, &timeslot.UpdatedAt)
	if err != nil {
		r.logger.Log("error while scanning")
		return nil, err
	}
	return &timeslot, nil
}
func (r repo) CheckAvailableTimeSlot(ctx context.Context, checkExist api.CheckExistedTimeslotRequest) (bool, error) {
	query := `SELECT id, start_time, booked_date, status, barber_id, created_at, updated_at FROM timeslot WHERE 
	  barber_id = $1
	  AND booked_date = $2
	  AND start_time = $3`

	var timeslot entity.Timeslot
	err := r.db.QueryRow(query, checkExist.BarberId, checkExist.BookedDate, checkExist.StartTime).
		Scan(&timeslot.ID, &timeslot.StartTime, &timeslot.BookedDate, &timeslot.Status, &timeslot.BarberId, &timeslot.CreatedAt, &timeslot.UpdatedAt)
	if err != nil {
		r.logger.Log("error while querying")
		return false, err
	}
	return true, nil
}
func (r repo) UpdateTimeSlot(ctx context.Context, timeslot api.CreateOrUpdateTimeslotRequest) (entity.Timeslot, error) {
	updatedTime := time.Now().Unix()
	query := `UPDATE timeslot SET start_time = $1, booked_date = $2, status = $3, barber_id = $4, updated_at = $5 WHERE id = $6`
	_, err := r.db.Exec(query, timeslot.StartTime, timeslot.BookedDate, timeslot.Status, timeslot.BarberId, updatedTime, timeslot.ID)
	if err != nil {
		r.logger.Log("error while updating data")
		return entity.Timeslot{}, err
	}
	return entity.Timeslot{
		ID:         timeslot.ID,
		StartTime:  timeslot.StartTime,
		BookedDate: timeslot.BookedDate,
		Status:     timeslot.Status,
		BarberId:   timeslot.BarberId,
		UpdatedAt:  updatedTime,
	}, nil
}
