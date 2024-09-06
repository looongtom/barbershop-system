package previewimage

import (
	"DoAn/pkg/previewimage/entity"
	"context"
)

type PreviewImageRepository interface {
	UploadImages(ctx context.Context, previewImage entity.PreviewImage) (*entity.PreviewImage, error)
	GetPreviewImage(ctx context.Context, id int) (*entity.PreviewImage, error)
	GetListPreviewImageByAccountId(ctx context.Context, accountId int) ([]entity.PreviewImage, error)
}
