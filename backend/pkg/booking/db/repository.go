package db

import (
	"DoAn"
	"DoAn/api"
	"DoAn/entity"
	"DoAn/mapper"
	"context"
	"database/sql"
	"github.com/lib/pq"
	"time"

	"github.com/go-kit/kit/log"
	_ "github.com/lib/pq"
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
func (r *repo) CreateBooking(ctx context.Context, booking api.BookingRequest) (entity.Booking, error) {
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

func (r repo) GetTotalCountBooking(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM booking;`
	var total int
	err := r.db.QueryRow(query).Scan(&total)
	if err != nil {
		r.logger.Log("error while scanning")
		return 0, err
	}
	return total, nil
}

func (r repo) GetTotalCountBookingByUserId(ctx context.Context, id int) (int, error) {
	query := `SELECT COUNT(*) FROM booking where customer_id = $1;`
	var total int
	err := r.db.QueryRow(query, id).Scan(&total)
	if err != nil {
		r.logger.Log("error while scanning")
		return 0, err
	}
	return total, nil
}

func (r repo) GetTotalCountBookingByBarberId(ctx context.Context, id int) (int, error) {
	query := `SELECT COUNT(*) FROM booking where barber_id = $1;`
	var total int
	err := r.db.QueryRow(query, id).Scan(&total)
	if err != nil {
		r.logger.Log("error while scanning")
		return 0, err
	}
	return total, nil
}

func (r repo) GetListBooking(ctx context.Context, page, pageSize int) ([]mapper.BookingMapper, error) {
	offset := (page - 1) * pageSize
	var listBooking []mapper.BookingMapper
	rows, err := r.db.Query(`SELECT 
			b.*, 
			COALESCE(ARRAY_AGG(bd.service_id), '{}') AS list_services
		FROM 
			booking b
		LEFT JOIN 
			booking_detail bd ON b.id = bd.booking_id
		GROUP BY 
			b.id
		LIMIT $1 OFFSET $2;
		`, pageSize, offset)
	if err != nil {
		r.logger.Log("error while querying")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var booking mapper.BookingMapper
		err = rows.Scan(&booking.ID, &booking.CustomerID, &booking.BarberId, &booking.ResultId, &booking.Status, &booking.Price, &booking.SlotId, &booking.FeedBackId, &booking.CreatedAt, &booking.UpdatedAt, pq.Array(&booking.ListServices))
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

func (r repo) UpdateBookingDetailService(ctx context.Context, listService []int, bookingId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Log("error while beginning transaction")
		return err
	}

	_, err = tx.Exec("DELETE FROM booking_detail WHERE booking_id = $1", bookingId)
	if err != nil {
		tx.Rollback()
		r.logger.Log("error while deleting data")
		return err
	}

	query := "INSERT INTO booking_detail(booking_id, service_id) VALUES($1, $2)"
	for _, serviceId := range listService {
		_, err := tx.Exec(query, bookingId, serviceId)
		if err != nil {
			tx.Rollback()
			r.logger.Log("error while inserting data")
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		r.logger.Log("error while committing transaction")
		return err
	}

	return nil
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

func (r repo) FindBookingByBarber(ctx context.Context, findReq api.FindListBookingRequest) ([]mapper.BookingMapper, error) {
	offset := (findReq.Page - 1) * findReq.PageSize
	var listBooking []mapper.BookingMapper
	rows, err := r.db.Query(`SELECT 
			b.*, 
			COALESCE(ARRAY_AGG(bd.service_id), '{}') AS list_services
		FROM 
			booking b
		LEFT JOIN 
			booking_detail bd ON b.id = bd.booking_id
		WHERE b.barber_id  = $1
		GROUP BY 
			b.id
		ORDER BY b.updated_at DESC
		LIMIT $2 OFFSET $3 ;
		`, findReq.Account, findReq.PageSize, offset)
	if err != nil {
		r.logger.Log("error while querying")
		return nil, err
	}
	for rows.Next() {
		var booking mapper.BookingMapper
		err = rows.Scan(&booking.ID, &booking.CustomerID, &booking.BarberId, &booking.ResultId, &booking.Status, &booking.Price, &booking.SlotId, &booking.FeedBackId, &booking.CreatedAt, &booking.UpdatedAt, pq.Array(&booking.ListServices))
		if err != nil {
			r.logger.Log("error while scanning")
			return nil, err
		}
		listBooking = append(listBooking, booking)
	}
	return listBooking, nil
}

func (r repo) FindBookingByUser(ctx context.Context, findReq api.FindListBookingRequest) ([]mapper.BookingMapper, error) {
	offset := (findReq.Page - 1) * findReq.PageSize
	var listBooking []mapper.BookingMapper
	rows, err := r.db.Query(`SELECT 
			b.*, 
			COALESCE(ARRAY_AGG(bd.service_id), '{}') AS list_services
		FROM 
			booking b
		LEFT JOIN 
			booking_detail bd ON b.id = bd.booking_id
		WHERE b.customer_id = $1
		GROUP BY 
			b.id
		ORDER BY b.updated_at DESC
		LIMIT $2 OFFSET $3 ;
		`, findReq.Account, findReq.PageSize, offset)
	if err != nil {
		r.logger.Log("error while querying")
		return nil, err
	}
	for rows.Next() {
		var booking mapper.BookingMapper
		err = rows.Scan(&booking.ID, &booking.CustomerID, &booking.BarberId, &booking.ResultId, &booking.Status, &booking.Price, &booking.SlotId, &booking.FeedBackId, &booking.CreatedAt, &booking.UpdatedAt, pq.Array(&booking.ListServices))
		if err != nil {
			r.logger.Log("error while scanning")
			return nil, err
		}
		listBooking = append(listBooking, booking)
	}
	return listBooking, nil
}
