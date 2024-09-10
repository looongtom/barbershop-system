package service

import (
	"DoAn/pkg/previewimage"
	"DoAn/pkg/previewimage/api"
	"DoAn/pkg/previewimage/common"
	"DoAn/pkg/previewimage/entity"
	"DoAn/pkg/previewimage/pb"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"time"

	"github.com/go-kit/kit/log"
)

type PreviewImageService struct {
	repository  previewimage.PreviewImageRepository
	connAccount *grpc.ClientConn
	logger      log.Logger
}

func NewService(repo previewimage.PreviewImageRepository, logger log.Logger,
	conn *grpc.ClientConn) previewimage.PreviewImageService {
	return &PreviewImageService{
		repository:  repo,
		logger:      logger,
		connAccount: conn,
	}
}

func (p PreviewImageService) UploadImages(ctx context.Context, request api.UpdateImageRequest) (interface{}, error) {
	clientAccount := pb.NewUserServiceClient(p.connAccount)
	checkBarber, err := clientAccount.CheckExistedBarber(ctx, &pb.CheckExistedBarberRequest{Id: int32(request.AccountId)})
	if err != nil {
		fmt.Printf("error when checking account: %v", err)
		return nil, err
	}
	if checkBarber == nil || checkBarber.Value == common.RoleUnknown {
		return nil, errors.New("account does not exist")
	}

	//ts := time.Now()
	SelfImg, err := common.UploadImageToCloud(request.SelfImg)
	if err != nil {
		return nil, err
	}
	ShapeImg, err := common.UploadImageToCloud(request.ShapeImg)
	if err != nil {
		return nil, err
	}
	ColorImg, err := common.UploadImageToCloud(request.ColorImg)
	if err != nil {
		return nil, err
	}

	return struct {
		SelfImg  string `json:"self_img"`
		ShapeImg string `json:"shape_img"`
		ColorImg string `json:"color_img"`
	}{
		SelfImg:  SelfImg,
		ShapeImg: ShapeImg,
		ColorImg: ColorImg,
	}, nil
}
func (p PreviewImageService) CreatePreviewImage(ctx context.Context, request api.CreatePreviewImageRequest) (interface{}, error) {
	//call grpc api
	clientAccount := pb.NewUserServiceClient(p.connAccount)
	checkBarber, err := clientAccount.CheckExistedBarber(ctx, &pb.CheckExistedBarberRequest{Id: int32(request.AccountId)})
	if err != nil {
		fmt.Printf("error when checking account: %v", err)
		return nil, err
	}
	if checkBarber == nil {
		return nil, errors.New("account does not exist")
	}

	ts := time.Now()
	urlImage, err := common.UploadImageToCloud(request.Url)
	if err != nil {
		return nil, err
	}
	resp, err := p.repository.UploadImages(ctx, entity.PreviewImage{
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
