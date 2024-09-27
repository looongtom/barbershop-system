package servicing

import (
	"DoAn/pkg/servicing/api"
	"DoAn/pkg/servicing/entity"
	"context"
)

type ServicingRepository interface {
	CreateCategory(ctx context.Context, category string) (*entity.Category, error)
	CreateService(ctx context.Context, category entity.Servicing) error
	UpdateService(ctx context.Context, category entity.Servicing) (*entity.Servicing, error)
	UpdateCategory(ctx context.Context, category entity.Category) (*entity.Category, error)

	GetService(ctx context.Context, id string) (*entity.Servicing, error)
	GetCategory(ctx context.Context, id string) (*entity.Servicing, error)
	GetListCategory(ctx context.Context) ([]entity.Category, error)
	GetListService(ctx context.Context) ([]entity.Servicing, error)
	GetListServiceAndCategory(ctx context.Context) ([]api.GetListServiceAndCategoryResponse, error)
	//GetListServiceByBookingId(ctx context.Context, bookingId int) ([]entity.Servicing, error)
}
