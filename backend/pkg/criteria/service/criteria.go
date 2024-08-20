package service

import (
	"DoAn/entity"
	"DoAn/pkg/criteria"
	"DoAn/pkg/criteria/api"
	"DoAn/pkg/criteria/common"
	"context"
	"github.com/go-kit/kit/log"
	"time"
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

func (c CriteriaStruct) CreateOrUpdateCriteria(ctx context.Context, criteria api.CreateOrUpdateCriteria) (interface{}, error) {
	ts := time.Now()
	urlImage, err := common.UploadImageToCloud(criteria.Img)
	if err != nil {
		return nil, err
	}
	resp, err := c.repository.CreateCriteria(ctx, entity.Criteria{
		Name:       criteria.Name,
		Img:        urlImage,
		CetegoryId: criteria.CategoryId,
		CreatedAt:  ts.Unix(),
		UpdatedAt:  ts.Unix(),
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
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
