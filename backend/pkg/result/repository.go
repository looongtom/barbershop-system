package result

import (
	"DoAn/entity"
	"context"
)

type ResultRepository interface {
	CreateOrUpdateResult(ctx context.Context, result entity.Result) (*entity.Result, error)
	UpdateResult(ctx context.Context, result entity.Result) (*entity.Result, error)

	GetResultByBarberId(ctx context.Context, barberId int) (*entity.Result, error)
	GetResultByUserId(ctx context.Context, userId int) (*entity.Result, error)
	GetResultByBookingId(ctx context.Context, bookingId int) (*entity.Result, error)
	GetResultById(ctx context.Context, bookingId int) (*entity.Result, error)

	CreateOrUpdateImageResult(ctx context.Context, imageResult entity.ImageResult) (*entity.ImageResult, error)
	GetImageResultByResultId(ctx context.Context, resultId int) ([]entity.ImageResult, error)
	DeleteImageResultByResultId(ctx context.Context, resultId int) error
}
