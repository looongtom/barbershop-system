package endpoint

import (
	"context"

	"DoAn/api"
	"DoAn/entity"
	"DoAn/pb"
)

func DecodeGetPreviewImageByUserRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetPreviewImageRequest)
	accountId := int(req.AccountId)
	return accountId, nil
}

func DecodeSavePreviewImageRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.SavePreviewImageRequest)
	return api.SaveGenerateRequest{
		SelfImg:      req.SelfImg,
		ShapeImg:     req.ShapeImg,
		ColorImg:     req.ColorImg,
		GeneratedImg: req.GeneratedImg,
		AccountId:    int(req.AccountId),
	}, nil
}

func EncodeGetPreviewImageByUserResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.([]entity.PreviewImage)
	// previewImageList := make([]entity.PreviewImage, 0)
	// for _, v := range resp.PreviewImages {
	// 	previewImageList = append(previewImageList, entity.PreviewImage{
	// 		ID:           int(v.Id),
	// 		SelfImg:      v.SelfImg,
	// 		ShapeImg:     v.ShapeImg,
	// 		ColorImg:     v.ColorImg,
	// 		GeneratedImg: v.GeneratedImg,
	// 		CreatedAt:    v.CreatedAt,
	// 	})
	// }
	return resp, nil
}

func EncodeSavePreviewImageResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*entity.PreviewImage)
	return &pb.PreviewImage{
		Id:           int32(resp.ID),
		SelfImg:      resp.SelfImg,
		ShapeImg:     resp.ShapeImg,
		ColorImg:     resp.ColorImg,
		GeneratedImg: resp.GeneratedImg,
		CreatedAt:    resp.CreatedAt,
	}, nil
}
