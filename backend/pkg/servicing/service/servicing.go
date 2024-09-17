package service

import (
	"DoAn/pkg/servicing"
	"DoAn/pkg/servicing/entity"
	"context"
	"errors"
	"strconv"

	"github.com/go-kit/kit/log"
)

type ServicingStruct struct {
	repository servicing.ServicingRepository
	logger     log.Logger
}

func NewService(repo servicing.ServicingRepository, logger log.Logger) servicing.ServicingService {
	return &ServicingStruct{
		repository: repo,
		logger:     logger,
	}
}

func (s ServicingStruct) GetServicing(ctx context.Context, id string) (interface{}, error) {
	resp, err := s.repository.GetService(ctx, id)
	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}
	return resp, nil
}

func (s ServicingStruct) GetListServicing(ctx context.Context) (interface{}, error) {
	resp, err := s.repository.GetListService(ctx)
	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}
	return resp, nil
}

func (s ServicingStruct) CreateServicing(ctx context.Context, service entity.Servicing) (interface{}, error) {
	cate, err := s.repository.GetCategory(ctx, strconv.Itoa(service.CategoryID))
	if err != nil {
		return nil, err
	}
	if cate == nil {
		return nil, errors.New("Category not found")
	}

	if err := s.repository.CreateService(ctx, service); err != nil {
		errMsg := err.Error()
		return errMsg, err
	}
	return service, nil
}

func (s ServicingStruct) UpdateServicing(ctx context.Context, service entity.Servicing) (interface{}, error) {
	id := strconv.Itoa(service.ID)
	if id == "" {
		return nil, errors.New("ID is required")
	}
	_, err := s.repository.GetService(ctx, id)
	if err != nil {
		return nil, err
	}
	cate, err := s.repository.GetCategory(ctx, strconv.Itoa(service.CategoryID))
	if err != nil {
		return nil, err
	}
	if cate == nil {
		return nil, errors.New("Category not found")
	}
	resp, err := s.repository.UpdateService(ctx, service)
	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}
	return resp, nil
}

func (s ServicingStruct) DeleteServicing(ctx context.Context, id string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (s ServicingStruct) GetListServicingByCategory(ctx context.Context, categoryId int) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (s ServicingStruct) GetListCategory(ctx context.Context) (interface{}, error) {
	resp, err := s.repository.GetListCategory(ctx)
	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}
	return resp, nil
}
func (s ServicingStruct) GetCategory(ctx context.Context, id string) (interface{}, error) {
	resp, err := s.repository.GetCategory(ctx, id)
	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}
	return resp, nil
}

func (s ServicingStruct) CreateCategory(ctx context.Context, category string) (interface{}, error) {
	resp, err := s.repository.CreateCategory(ctx, category)
	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}
	return resp, nil
}

func (s ServicingStruct) UpdateCategory(ctx context.Context, category entity.Category) (interface{}, error) {
	id := strconv.Itoa(category.ID)
	if id == "" {
		return nil, errors.New("ID is required")
	}
	_, err := s.repository.GetCategory(ctx, id)
	if err != nil {
		return nil, err
	}
	resp, err := s.repository.UpdateCategory(ctx, category)
	if err != nil {
		errMsg := err.Error()
		return errMsg, err
	}
	return resp, nil
}

func (s ServicingStruct) DeleteCategory(ctx context.Context, id string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}
