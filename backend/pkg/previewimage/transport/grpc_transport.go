package transport

import (
	"context"

	"DoAn/endpoint"
	"DoAn/entity"

	gt "github.com/go-kit/kit/transport/grpc"

	"DoAn/pb"
)

type GRPCServer struct {
	getPreviewImageByUser gt.Handler
	savePreviewImage      gt.Handler
}

func (G GRPCServer) GetPreviewImageByUser(ctx context.Context, request *pb.GetPreviewImageRequest) (*pb.PreviewImageList, error) {
	_, resp, err := G.getPreviewImageByUser.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	var listGrpcPreviewImage []*pb.PreviewImage
	for _, v := range resp.([]entity.PreviewImage) {
		listGrpcPreviewImage = append(listGrpcPreviewImage, &pb.PreviewImage{
			Id:           int32(v.ID),
			SelfImg:      v.SelfImg,
			ShapeImg:     v.ShapeImg,
			ColorImg:     v.ColorImg,
			GeneratedImg: v.GeneratedImg,
			CreatedAt:    v.CreatedAt,
		})
	}
	return &pb.PreviewImageList{PreviewImages: listGrpcPreviewImage}, nil
}

func (G GRPCServer) SavePreviewImage(ctx context.Context, request *pb.SavePreviewImageRequest) (*pb.PreviewImage, error) {
	_, resp, err := G.savePreviewImage.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.PreviewImage), nil
}

func NewGRPCServer(_ context.Context, endpoints endpoint.Endpoints) pb.PreviewImageServiceServer {
	return &GRPCServer{
		getPreviewImageByUser: gt.NewServer(
			endpoints.GetPreviewImageByUserEndpoint,
			endpoint.DecodeGetPreviewImageByUserRequest,
			endpoint.EncodeGetPreviewImageByUserResponse,
		),
		savePreviewImage: gt.NewServer(
			endpoints.SavePreviewImageEndpoint,
			endpoint.DecodeSavePreviewImageRequest,
			endpoint.EncodeSavePreviewImageResponse,
		),
	}
}
