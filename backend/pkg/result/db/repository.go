package db

import (
	"DoAn"
	"DoAn/entity"
	"context"
	"database/sql"
	"time"

	"github.com/go-kit/kit/log"
)

type repo struct {
	db     *sql.DB
	logger log.Logger
}

func NewRepository(db *sql.DB, logger log.Logger) (result.ResultRepository, error) {
	return &repo{
		db:     db,
		logger: logger,
	}, nil
}

func (r repo) CreateOrUpdateResult(ctx context.Context, result entity.Result) (*entity.Result, error) {
	query := `INSERT INTO result (barber_id, user_id, booking_id, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, result.BarberId, result.UserId, result.BookingId, result.Description, result.CreatedAt, result.UpdatedAt).Scan(&result.Id)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r repo) GetResultByBarberId(ctx context.Context, barberId int) (entity.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) GetResultByUserId(ctx context.Context, userId int) (entity.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) GetResultByBookingId(ctx context.Context, bookingId int) (*entity.Result, error) {
	query := `SELECT * FROM result WHERE booking_id=$1`
	var result entity.Result
	err := r.db.QueryRowContext(ctx, query, bookingId).Scan(&result.Id, &result.BarberId, &result.UserId, &result.BookingId, &result.Description, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r repo) GetImageResultByResultId(ctx context.Context, resultId int) ([]entity.ImageResult, error) {
	query := `SELECT * FROM image_result WHERE result_id=$1`
	rows, err := r.db.QueryContext(ctx, query, resultId)
	if err != nil {
		return nil, err
	}
	var listImageResult []entity.ImageResult
	for rows.Next() {
		var imageResult entity.ImageResult
		err = rows.Scan(&imageResult.Id, &imageResult.Url, &imageResult.ResultId, &imageResult.CreatedAt)
		if err != nil {
			return nil, err
		}
		listImageResult = append(listImageResult, imageResult)
	}
	return listImageResult, nil
}
func (r repo) CreateOrUpdateImageResult(ctx context.Context, imageResult entity.ImageResult) (*entity.ImageResult, error) {
	query := `INSERT INTO image_result (url, result_id, created_at) VALUES ($1, $2, $3) RETURNING id`
	ts := time.Now()
	err := r.db.QueryRowContext(ctx, query, imageResult.Url, imageResult.ResultId, ts.Unix()).Scan(&imageResult.Id)
	if err != nil {
		return nil, err
	}
	return &imageResult, nil
}
