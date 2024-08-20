package criteria

import (
	"DoAn/entity"
	"DoAn/pkg/criteria/api"
	"context"
)

type CriteriaService interface {
	GetCriteria(ctx context.Context, id string) (interface{}, error)
	GetListCriteria(ctx context.Context) (interface{}, error)
	CreateOrUpdateCriteria(ctx context.Context, criteria entity.Criteria) (interface{}, error)
	DeleteCriteria(ctx context.Context, id string) (interface{}, error)
	GetListCriteriaByCategory(ctx context.Context, categoryId int) (interface{}, error)

	GetCategory(ctx context.Context, id string) (interface{}, error)
	GetListCategory(ctx context.Context) (interface{}, error)
	CreateOrUpdateCategory(ctx context.Context, category api.CreateOrUpdateCategory) (interface{}, error)
	DeleteCategory(ctx context.Context, id string) (interface{}, error)
}
