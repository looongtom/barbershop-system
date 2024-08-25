package criteria

import (
	"DoAn/pkg/criteria/api"
	"context"
)

type CriteriaService interface {
	GetCriteria(ctx context.Context, id string) (interface{}, error)
	GetListCriteria(ctx context.Context) (interface{}, error)
	CreateOrUpdateCriteria(ctx context.Context, criteria api.CreateOrUpdateCriteria) (interface{}, error)
	DeleteCriteria(ctx context.Context, id string) (interface{}, error)
	FindCriteria(ctx context.Context, findCriteriaReq api.FindCriteria) (interface{}, error)

	GetCategory(ctx context.Context, id string) (interface{}, error)
	GetListCategory(ctx context.Context) (interface{}, error)
	CreateOrUpdateCategory(ctx context.Context, category api.CreateOrUpdateCategory) (interface{}, error)
	DeleteCategory(ctx context.Context, id string) (interface{}, error)
}
