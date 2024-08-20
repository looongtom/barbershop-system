package service

import (
	"DoAn/pkg/result"
	"DoAn/pkg/result/api"
	"context"

	"github.com/go-kit/kit/log"
)

type ResultService struct {
	repository result.ResultRepository
	logger     log.Logger
}

func NewService(repo result.ResultRepository, logger log.Logger) result.ResultService {
	return &ResultService{
		repository: repo,
		logger:     logger,
	}
}

func (r ResultService) CreateOrUpdateResult(ctx context.Context, result api.CreateOrUpdateResult) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (r ResultService) GetResultByBarberId(ctx context.Context, barberId int) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (r ResultService) GetResultByUserId(ctx context.Context, userId int) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (r ResultService) GetResultByBookingId(ctx context.Context, bookingId int) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (r ResultService) CreateOrUpdateImageResult(ctx context.Context, imageResult api.CreateImageResult) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}
