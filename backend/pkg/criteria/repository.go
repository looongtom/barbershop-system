package criteria

import (
	entity2 "DoAn/pkg/criteria/entity"
	"DoAn/pkg/servicing/entity"
	"context"
)

type CriteriaRepository interface {
	CreateCategory(ctx context.Context, category string) (*entity.Category, error)
	CreateCriteria(ctx context.Context, criteria entity2.Criteria) (*entity2.Criteria, error)
	UpdateCriteria(ctx context.Context, criteria entity2.Criteria) (*entity2.Criteria, error)
	UpdateCategory(ctx context.Context, category entity.Category) (*entity.Category, error)

	GetCriteria(ctx context.Context, id string) (*entity2.Criteria, error)
	GetCategory(ctx context.Context, id string) (*entity.Category, error)
	GetListCategory(ctx context.Context) ([]entity.Category, error)
	GetListCriteria(ctx context.Context) ([]entity2.Criteria, error)

	FindCriteria(ctx context.Context, name string, category int) ([]entity2.Criteria, error)
}
