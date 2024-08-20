package criteria

import (
	"DoAn/entity"
	"context"
)

type CriteriaRepository interface {
	CreateCategory(ctx context.Context, category string) (*entity.Category, error)
	CreateCriteria(ctx context.Context, criteria entity.Criteria) (*entity.Criteria, error)
	UpdateCriteria(ctx context.Context, criteria entity.Criteria) (*entity.Criteria, error)
	UpdateCategory(ctx context.Context, category entity.Category) (*entity.Category, error)

	GetCriteria(ctx context.Context, id string) (*entity.Criteria, error)
	GetCategory(ctx context.Context, id string) (*entity.Criteria, error)
	GetListCategory(ctx context.Context) ([]entity.Category, error)
	GetListCriteria(ctx context.Context) ([]entity.Criteria, error)
}
