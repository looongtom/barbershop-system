package previewimage

import (
	"DoAn/pkg/previewimage/api"
	"context"
)

type PreviewImageService interface {
	CreatePreviewImage(ctx context.Context, request api.CreatePreviewImageRequest) (interface{}, error)
	UploadImages(ctx context.Context, request api.UpdateImageRequest) (interface{}, error)
	GetPreviewImage(ctx context.Context, id int) (interface{}, error)
	GetListPreviewImageByAccountId(ctx context.Context, accountId int) (interface{}, error)
}
