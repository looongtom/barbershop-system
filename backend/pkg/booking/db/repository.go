package db

import (
	"DoAn/pkg/booking"
	"DoAn/pkg/booking/api"
	"DoAn/pkg/booking/entity"
	"context"
	"database/sql"
	"time"

	"github.com/go-kit/kit/log"
)

type repo struct {
	db     *sql.DB
	logger log.Logger
}

func NewRepository(db *sql.DB, logger log.Logger) (booking.BookingRepository, error) {
	return &repo{
		db:     db,
		logger: logger,
	}, nil
}

func (r repo) GetBookingDetail(ctx context.Context, bookingId int) ([]entity.BookingDetail, error) {
	var listBookingDetail []entity.BookingDetail
	rows, err := r.db.Query("SELECT * FROM booking_detail WHERE booking_id = $1", bookingId)
	if err != nil {
		r.logger.Log("error while querying")
		return nil, err
	}
	for rows.Next() {
		var bookingDetail entity.BookingDetail
		err = rows.Scan(&bookingDetail.BookingId, &bookingDetail.ServiceId)
		if err != nil {
			r.logger.Log("error while scanning")
			return nil, err
		}
		listBookingDetail = append(listBookingDetail, bookingDetail)
	}
	return listBookingDetail, nil

}
func (r repo) CreateBookingDetail(ctx context.Context, listService []int, bookingId int) error {
	query := `
	INSERT INTO booking_detail(booking_id, service_id) VALUES($1, $2)
`
	for _, serviceId := range listService {
		_, err := r.db.Exec(query, bookingId, serviceId)
		if err != nil {
			r.logger.Log("error while inserting data")
			return err
		}
	}
	return nil
}
func (r repo) CreateBooking(ctx context.Context, booking api.BookingRequest) (entity.Booking, error) {
	newBooking := entity.Booking{
		CustomerID: booking.CustomerID,
		BarberId:   booking.BarberId,
		ResultId:   booking.ResultId,
		Status:     booking.Status,
		Price:      booking.Price,
		SlotId:     booking.SlotId,
		FeedBackId: booking.FeedBackId,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}
	_, err := r.db.Exec("INSERT INTO booking(customer_id, barber_id, result_id, status, price, slot_id, feedback_id, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)", newBooking.CustomerID, newBooking.BarberId, newBooking.ResultId, newBooking.Status, newBooking.Price, newBooking.SlotId, newBooking.FeedBackId, newBooking.CreatedAt, newBooking.UpdatedAt)
	if err != nil {
		r.logger.Log("error while inserting data")
		return entity.Booking{}, err
	}
	//get booking id after insert
	var bookingId int
	err = r.db.QueryRow("SELECT id FROM booking WHERE created_at = $1", newBooking.CreatedAt).Scan(&bookingId)
	if err != nil {
		r.logger.Log("error while scanning")
		return entity.Booking{}, err
	}
	newBooking.ID = bookingId
	return newBooking, nil
}

func (r repo) GetListIdServiceByBookingId(ctx context.Context, id int) ([]int, error) {
	var listIds []int
	rows, err := r.db.Query("SELECT service_id FROM booking_detail WHERE booking_id = $1", id)
	if err != nil {
		r.logger.Log("error while querying")
		return nil, err
	}
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			r.logger.Log("error while scanning")
			return nil, err
		}
		listIds = append(listIds, id)
	}
	return listIds, nil
}
func (r repo) GetListBooking(ctx context.Context) ([]entity.Booking, error) {
	var listBooking []entity.Booking
	rows, err := r.db.Query("SELECT * FROM booking")
	if err != nil {
		r.logger.Log("error while querying")
		return nil, err
	}
	for rows.Next() {
		var booking entity.Booking
		err = rows.Scan(&booking.ID, &booking.CustomerID, &booking.BarberId, &booking.ResultId, &booking.Status, &booking.Price, &booking.SlotId, &booking.FeedBackId, &booking.CreatedAt, &booking.UpdatedAt)
		if err != nil {
			r.logger.Log("error while scanning")
			return nil, err
		}
		listBooking = append(listBooking, booking)
	}
	return listBooking, nil
}

func (r repo) GetBookingById(ctx context.Context, id int) (entity.Booking, error) {
	var booking entity.Booking
	err := r.db.QueryRow("SELECT * FROM booking WHERE id = $1", id).Scan(&booking.ID, &booking.CustomerID, &booking.BarberId, &booking.ResultId, &booking.Status, &booking.Price, &booking.SlotId, &booking.FeedBackId, &booking.CreatedAt, &booking.UpdatedAt)
	if err != nil {
		r.logger.Log("error while scanning")
		return entity.Booking{}, err
	}
	return booking, nil
}

func (r repo) UpdateBooking(ctx context.Context, booking api.UpdateBookingRequest) (entity.Booking, error) {
	query := `UPDATE booking SET customer_id = $1, barber_id = $2, result_id = $3, status = $4, price = $5, slot_id = $6, feedback_id = $7, updated_at = $8 WHERE id = $9`
	updateBooking := entity.Booking{
		ID:         booking.Id,
		CustomerID: booking.CustomerID,
		BarberId:   booking.BarberId,
		ResultId:   booking.ResultId,
		Status:     booking.Status,
		Price:      booking.Price,
		SlotId:     booking.SlotId,
		FeedBackId: booking.FeedBackId,
		UpdatedAt:  time.Now().Unix(),
	}
	_, err := r.db.Exec(query, updateBooking.CustomerID, updateBooking.BarberId, updateBooking.ResultId, updateBooking.Status, updateBooking.Price, updateBooking.SlotId, updateBooking.FeedBackId, updateBooking.UpdatedAt, updateBooking.ID)
	if err != nil {
		r.logger.Log("error while updating data")
		return entity.Booking{}, err
	}
	return updateBooking, nil
}

func (r repo) FindBookingByUserOrBarber(ctx context.Context, findReq api.FindBookingRequest) ([]entity.Booking, error) {
	var listBooking []entity.Booking
	rows, err := r.db.Query("SELECT * FROM booking WHERE customer_id = $1 OR barber_id = $2", findReq.CustomerId, findReq.BarberId)
	if err != nil {
		r.logger.Log("error while querying")
		return nil, err
	}
	for rows.Next() {
		var booking entity.Booking
		err = rows.Scan(&booking.ID, &booking.CustomerID, &booking.BarberId, &booking.ResultId, &booking.Status, &booking.Price, &booking.SlotId, &booking.FeedBackId, &booking.CreatedAt, &booking.UpdatedAt)
		if err != nil {
			r.logger.Log("error while scanning")
			return nil, err
		}
		listBooking = append(listBooking, booking)
	}
	return listBooking, nil
}
