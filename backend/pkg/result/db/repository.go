package db

import (
	"DoAn"
	"DoAn/api"
	"DoAn/entity"
	"context"
	"database/sql"

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

func (r repo) CreateOrUpdateResult(ctx context.Context, result entity.Result) (entity.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) GetResultByBarberId(ctx context.Context, barberId int) (entity.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) GetResultByUserId(ctx context.Context, userId int) (entity.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) GetResultByBookingId(ctx context.Context, bookingId int) (entity.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) CreateOrUpdateImageResult(ctx context.Context, imageResult api.CreateImageResult) (entity.ImageResult, error) {
	//TODO implement me
	panic("implement me")
}
