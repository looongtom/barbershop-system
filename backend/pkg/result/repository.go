package result

import (
	"DoAn/api"
	"DoAn/entity"
	"context"
)

type ResultRepository interface {
	CreateOrUpdateResult(ctx context.Context, result entity.Result) (entity.Result, error)

	GetResultByBarberId(ctx context.Context, barberId int) (entity.Result, error)
	GetResultByUserId(ctx context.Context, userId int) (entity.Result, error)
	GetResultByBookingId(ctx context.Context, bookingId int) (entity.Result, error)

	CreateOrUpdateImageResult(ctx context.Context, imageResult api.CreateImageResult) (entity.ImageResult, error)
}
