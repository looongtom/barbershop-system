package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	previewimage "DoAn"
	"DoAn/api"
)

type Endpoints struct {
	GetPreviewImageByUserEndpoint endpoint.Endpoint
	SavePreviewImageEndpoint      endpoint.Endpoint
}

func MakeGetPreviewImageByUserEndpoint(svc previewimage.PreviewImageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(int)
		return svc.GetListPreviewImageByAccountId(ctx, req)
	}
}

func MakeSavePreviewImageEndpoint(svc previewimage.PreviewImageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.SaveGenerateRequest)
		return svc.SaveGeneratedImage(ctx, req)
	}
}
