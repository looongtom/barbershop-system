package service

import (
	"DoAn"
	"DoAn/api"
	"DoAn/common"
	"DoAn/entity"
	"context"
	"time"

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

func (r ResultService) UpdateResult(ctx context.Context, result api.CreateOrUpdateResult) (interface{}, error) {
	ts := time.Now()
	_, err := r.repository.GetResultById(ctx, result.ID)
	if err != nil {
		return nil, err
	}
	r.repository.UpdateResult(ctx, entity.Result{
		Id:          result.ID,
		Description: result.Description,
		UpdatedAt:   ts.Unix(),
	})

	var listUrlImg []string
	if result.ListImg != nil {

		err = r.repository.DeleteImageResultByResultId(ctx, result.ID)
		if err != nil {
			return nil, err
		}

		for _, img := range result.ListImg {
			urlImage, err := common.UploadImageToCloud(img)
			if err != nil {
				r.logger.Log("error while uploading image")
				return nil, err
			}
			listUrlImg = append(listUrlImg, urlImage)
		}

		var listImgResultResp []api.CreateOrUpdateImageResponse
		for _, url := range listUrlImg {
			imgResult, err := r.repository.CreateOrUpdateImageResult(ctx, entity.ImageResult{
				ResultId:  result.ID,
				Url:       url,
				CreatedAt: ts.Unix(),
			})
			if err != nil {
				return nil, err
			}
			listImgResultResp = append(listImgResultResp, api.CreateOrUpdateImageResponse{
				ID:       imgResult.Id,
				Url:      imgResult.Url,
				ResultId: imgResult.ResultId,
			})
		}
		return api.CreateOrUpdateResultResponse{
			ID:          result.ID,
			BarberId:    result.BarberId,
			UserId:      result.UserId,
			BookingId:   result.BookingId,
			Description: result.Description,
			ListImg:     listImgResultResp,
		}, nil
	}
	return api.CreateOrUpdateResultResponse{
		ID:          result.ID,
		BarberId:    result.BarberId,
		UserId:      result.UserId,
		BookingId:   result.BookingId,
		Description: result.Description,
	}, nil
}

func (r ResultService) CreateResult(ctx context.Context, result api.CreateOrUpdateResult) (interface{}, error) {
	ts := time.Now()
	var listUrlImg []string
	for _, img := range result.ListImg {
		urlImage, err := common.UploadImageToCloud(img)
		if err != nil {
			r.logger.Log("error while uploading image")
			return nil, err
		}
		listUrlImg = append(listUrlImg, urlImage)
	}
	respResult, err := r.repository.CreateOrUpdateResult(ctx, entity.Result{
		Id:          0,
		BarberId:    result.BarberId,
		UserId:      result.UserId,
		BookingId:   result.BookingId,
		Description: result.Description,
		CreatedAt:   ts.Unix(),
		UpdatedAt:   ts.Unix(),
	})
	if err != nil {
		return nil, err
	}
	var listImgResultResp []api.CreateOrUpdateImageResponse
	for _, url := range listUrlImg {
		imgResult, err := r.repository.CreateOrUpdateImageResult(ctx, entity.ImageResult{
			ResultId:  respResult.Id,
			Url:       url,
			CreatedAt: ts.Unix(),
		})
		if err != nil {
			return nil, err
		}
		listImgResultResp = append(listImgResultResp, api.CreateOrUpdateImageResponse{
			ID:       imgResult.Id,
			Url:      imgResult.Url,
			ResultId: imgResult.ResultId,
		})
	}
	return api.CreateOrUpdateResultResponse{
		ID:          result.ID,
		BarberId:    result.BarberId,
		UserId:      result.UserId,
		BookingId:   result.BookingId,
		Description: result.Description,
		ListImg:     listImgResultResp,
	}, nil
}

func (r ResultService) GetResultByBarberId(ctx context.Context, barberId int) (interface{}, error) {
	resp, err := r.repository.GetResultByBarberId(ctx, barberId)
	if err != nil {
		return nil, err
	}
	listImgs, err := r.repository.GetImageResultByResultId(ctx, resp.Id)
	if err != nil {
		return nil, err
	}
	var listImgResultResp []api.GetResultResponse
	for _, img := range listImgs {
		listImgResultResp = append(listImgResultResp, api.GetResultResponse{
			ID:       img.Id,
			Url:      img.Url,
			ResultId: img.ResultId,
		})

	}
	return api.GetDetailResultResponse{
		ID:          resp.Id,
		BarberId:    resp.BarberId,
		UserId:      resp.UserId,
		BookingId:   resp.BookingId,
		Description: resp.Description,
		ListImg:     listImgResultResp,
	}, err

}

func (r ResultService) GetResultByUserId(ctx context.Context, userId int) (interface{}, error) {
	resp, err := r.repository.GetResultByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	listImgs, err := r.repository.GetImageResultByResultId(ctx, resp.Id)
	if err != nil {
		return nil, err
	}
	var listImgResultResp []api.GetResultResponse
	for _, img := range listImgs {
		listImgResultResp = append(listImgResultResp, api.GetResultResponse{
			ID:       img.Id,
			Url:      img.Url,
			ResultId: img.ResultId,
		})

	}
	return api.GetDetailResultResponse{
		ID:          resp.Id,
		BarberId:    resp.BarberId,
		UserId:      resp.UserId,
		BookingId:   resp.BookingId,
		Description: resp.Description,
		ListImg:     listImgResultResp,
	}, err

}

func (r ResultService) GetResultByBookingId(ctx context.Context, bookingId int) (interface{}, error) {
	resp, err := r.repository.GetResultByBookingId(ctx, bookingId)
	if err != nil {
		return nil, err
	}
	listImgs, err := r.repository.GetImageResultByResultId(ctx, resp.Id)
	if err != nil {
		return nil, err
	}
	var listImgResultResp []api.GetResultResponse
	for _, img := range listImgs {
		listImgResultResp = append(listImgResultResp, api.GetResultResponse{
			ID:       img.Id,
			Url:      img.Url,
			ResultId: img.ResultId,
		})

	}
	return api.GetDetailResultResponse{
		ID:          resp.Id,
		BarberId:    resp.BarberId,
		UserId:      resp.UserId,
		BookingId:   resp.BookingId,
		Description: resp.Description,
		ListImg:     listImgResultResp,
	}, err

}

func (r ResultService) CreateOrUpdateImageResult(ctx context.Context, imageResult api.CreateImageResult) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}
