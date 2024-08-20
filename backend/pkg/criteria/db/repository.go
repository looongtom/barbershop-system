package db

import (
	"DoAn/entity"
	"DoAn/pkg/criteria"
	"context"
	"database/sql"
	"github.com/go-kit/kit/log"
)

type repo struct {
	db     *sql.DB
	logger log.Logger
}

func NewRepository(db *sql.DB, logger log.Logger) (criteria.CriteriaRepository, error) {
	return &repo{
		db:     db,
		logger: logger,
	}, nil
}

func (r repo) CreateCategory(ctx context.Context, category string) (*entity.Category, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) CreateCriteria(ctx context.Context, criteria entity.Criteria) (*entity.Criteria, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) UpdateCriteria(ctx context.Context, criteria entity.Criteria) (*entity.Criteria, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) UpdateCategory(ctx context.Context, category entity.Category) (*entity.Category, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) GetCriteria(ctx context.Context, id string) (*entity.Criteria, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) GetCategory(ctx context.Context, id string) (*entity.Criteria, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) GetListCategory(ctx context.Context) ([]entity.Category, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) GetListCriteria(ctx context.Context) ([]entity.Criteria, error) {
	//TODO implement me
	panic("implement me")
}
