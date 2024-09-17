package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"DoAn/pkg/previewimage"
	"DoAn/pkg/previewimage/api"
	"DoAn/pkg/previewimage/common"
	"DoAn/pkg/previewimage/entity"
	kafka2 "DoAn/pkg/previewimage/kafka"
	"DoAn/pkg/previewimage/pb"

	"google.golang.org/grpc"

	"github.com/go-kit/kit/log"
)

const (
	topic = "preview_img"
)

type PreviewImageService struct {
	repository  previewimage.PreviewImageRepository
	connAccount *grpc.ClientConn
	logger      log.Logger
	kafka       *kafka.Producer
}

func NewService(repo previewimage.PreviewImageRepository, logger log.Logger,
	conn *grpc.ClientConn, kafka *kafka.Producer) previewimage.PreviewImageService {
	return &PreviewImageService{
		repository:  repo,
		logger:      logger,
		connAccount: conn,
		kafka:       kafka,
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

	// ts := time.Now()
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

	previewImg := entity.PreviewImage{
		CreatedAt: time.Now().Unix(),
		AccountId: request.AccountId,
		SelfImg:   SelfImg,
		ShapeImg:  ShapeImg,
		ColorImg:  ColorImg,
	}
	serializedPreviewImg, err := json.Marshal(previewImg)

	err = kafka2.ProduceMessage(p.kafka, topic, serializedPreviewImg)
	if err != nil {
		p.logger.Log("Failed to produce message: %s\n", err)
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
	// call grpc api
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
