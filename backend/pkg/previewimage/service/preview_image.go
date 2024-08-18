package service

import (
	"DoAn/entity"
	"DoAn/pkg/previewimage"
	"DoAn/pkg/previewimage/api"
	"DoAn/pkg/previewimage/common"
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type PreviewImageService struct {
	repository previewimage.PreviewImageRepository
	logger     log.Logger
}

func (p PreviewImageService) CreatePreviewImage(ctx context.Context, request api.CreatePreviewImageRequest) (interface{}, error) {
	ts := time.Now()
	urlImage, err := common.UploadImageToCloud(request.Url)
	if err != nil {
		return nil, err
	}
	resp, err := p.repository.CreatePreviewImage(ctx, entity.PreviewImage{
		Url:       urlImage,
		CreatedAt: ts.Unix(),
		AccountId: request.AccountId,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (p PreviewImageService) GetPreviewImage(ctx context.Context, id int) (interface{}, error) {
	resp, err := p.repository.GetPreviewImage(ctx, id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (p PreviewImageService) GetListPreviewImageByAccountId(ctx context.Context, accountId int) (interface{}, error) {
	resp, err := p.repository.GetListPreviewImageByAccountId(ctx, accountId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewService(repo previewimage.PreviewImageRepository, logger log.Logger) previewimage.PreviewImageService {
	return &PreviewImageService{
		repository: repo,
		logger:     logger,
	}
}
