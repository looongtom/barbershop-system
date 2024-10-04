package criteria

import (
	"DoAn/entity"
	"context"
)

type CriteriaRepository interface {
	CreateCategory(ctx context.Context, CategoryCriteria string) (*entity.CategoryCriteria, error)
	CreateCriteria(ctx context.Context, criteria entity.Criteria) (*entity.Criteria, error)
	UpdateCriteria(ctx context.Context, criteria entity.Criteria) (*entity.Criteria, error)
	UpdateCategory(ctx context.Context, CategoryCriteria entity.CategoryCriteria) (*entity.CategoryCriteria, error)

	GetCriteria(ctx context.Context, id string) (*entity.Criteria, error)
	GetCategory(ctx context.Context, id string) (*entity.CategoryCriteria, error)
	GetListCategory(ctx context.Context) ([]entity.CategoryCriteria, error)
	GetListCriteria(ctx context.Context) ([]entity.Criteria, error)

	FindCriteria(ctx context.Context, name string, CategoryCriteria int) ([]entity.Criteria, error)
}
