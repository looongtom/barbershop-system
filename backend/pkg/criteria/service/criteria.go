package service

import (
	"DoAn/entity"
	"DoAn/pkg/criteria"
	"DoAn/pkg/criteria/api"
	"context"
	"github.com/go-kit/kit/log"
)

type CriteriaStruct struct {
	repository criteria.CriteriaRepository
	logger     log.Logger
}

func NewService(repo criteria.CriteriaRepository, logger log.Logger) criteria.CriteriaService {
	return &CriteriaStruct{
		repository: repo,
		logger:     logger,
	}
}

func (c CriteriaStruct) GetCriteria(ctx context.Context, id string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (c CriteriaStruct) GetListCriteria(ctx context.Context) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (c CriteriaStruct) CreateOrUpdateCriteria(ctx context.Context, criteria entity.Criteria) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (c CriteriaStruct) DeleteCriteria(ctx context.Context, id string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (c CriteriaStruct) GetListCriteriaByCategory(ctx context.Context, categoryId int) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (c CriteriaStruct) GetCategory(ctx context.Context, id string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (c CriteriaStruct) GetListCategory(ctx context.Context) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (c CriteriaStruct) CreateOrUpdateCategory(ctx context.Context, category api.CreateOrUpdateCategory) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (c CriteriaStruct) DeleteCategory(ctx context.Context, id string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}
