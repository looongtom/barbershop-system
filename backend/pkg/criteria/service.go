package criteria

import (
	"DoAn/entity"
	"context"
)

type CriteriaService interface {
	GetCriteria(ctx context.Context, id string) (interface{}, error)
	GetListCriteria(ctx context.Context) (interface{}, error)
	CreateCriteria(ctx context.Context, criteria entity.Criteria) (interface{}, error)
	UpdateCriteria(ctx context.Context, criteria entity.Criteria) (interface{}, error)
	DeleteCriteria(ctx context.Context, id string) (interface{}, error)
	GetListCriteriaByCategory(ctx context.Context, categoryId int) (interface{}, error)

	GetCategory(ctx context.Context, id string) (interface{}, error)
	GetListCategory(ctx context.Context) (interface{}, error)
	CreateCategory(ctx context.Context, category string) (interface{}, error)
	UpdateCategory(ctx context.Context, category entity.Category) (interface{}, error)
	DeleteCategory(ctx context.Context, id string) (interface{}, error)
}
