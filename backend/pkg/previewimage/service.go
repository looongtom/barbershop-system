package previewimage

import (
	"DoAn/api"
	"context"
)

type PreviewImageService interface {
	//CreatePreviewImage(ctx context.Context, request api.CreatePreviewImageRequest) (interface{}, error)
	UploadImages(ctx context.Context, request api.UpdateImageRequest) (interface{}, error)
	GetPreviewImage(ctx context.Context, id int) (interface{}, error)
	GetListPreviewImageByAccountId(ctx context.Context, accountId int) (interface{}, error)
	SaveGeneratedImage(ctx context.Context, request api.SaveGenerateRequest) (interface{}, error)
}
