package servicing

import (
	"DoAn/entity"
	"context"
)

type ServicingService interface {
	GetServicing(ctx context.Context, id string) (interface{}, error)
	GetListServicing(ctx context.Context) (interface{}, error)
	CreateServicing(ctx context.Context, servicing entity.Servicing) (interface{}, error)
	UpdateServicing(ctx context.Context, servicing entity.Servicing) (interface{}, error)
	DeleteServicing(ctx context.Context, id string) (interface{}, error)
	GetListServicingByCategory(ctx context.Context, categoryId int) (interface{}, error)
	//GetListServicingByBookingId(ctx context.Context, bookingId int) (interface{}, error)
	GetListServicingAndCategory(ctx context.Context) (interface{}, error)

	GetCategory(ctx context.Context, id string) (interface{}, error)
	GetListCategory(ctx context.Context) (interface{}, error)
	CreateCategory(ctx context.Context, category string) (interface{}, error)
	UpdateCategory(ctx context.Context, category entity.Category) (interface{}, error)
	DeleteCategory(ctx context.Context, id string) (interface{}, error)
}
