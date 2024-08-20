package result

import (
	"DoAn/pkg/result/api"
	"context"
)

type ResultService interface {
	CreateOrUpdateResult(ctx context.Context, result api.CreateOrUpdateResult) (interface{}, error)
	GetResultByBarberId(ctx context.Context, barberId int) (interface{}, error)
	GetResultByUserId(ctx context.Context, userId int) (interface{}, error)
	GetResultByBookingId(ctx context.Context, bookingId int) (interface{}, error)

	CreateOrUpdateImageResult(ctx context.Context, imageResult api.CreateImageResult) (interface{}, error)
}
