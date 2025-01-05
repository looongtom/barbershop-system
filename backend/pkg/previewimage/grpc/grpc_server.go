package main

import (
	"context"
	"fmt"
	logV "log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	previewimage "DoAn"
	"DoAn/database"
	repository "DoAn/db"
	"DoAn/endpoint"
	"DoAn/pb"
	"DoAn/service"
	"DoAn/transport"

	"github.com/go-kit/kit/log"
)

func main() {
	err := godotenv.Load("previewimage.env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}
	logger := log.NewLogfmtLogger(os.Stderr)
	collectionPostgres, err := database.ConnectPostgresPreviewImage()
	if err != nil {
		logV.Fatalf("Error getting env, %v", err)
	}
	ctx := context.Background()

	var svc previewimage.PreviewImageService
	svc = service.PreviewImageService{}
	{
		repo, err := repository.NewRepository(collectionPostgres, logger)
		if err != nil {
			logV.Fatalf("Error getting env, %v", err)
		}
		svc = service.NewService(repo, logger, nil, nil)
	}
	errors := make(chan error)
	go func() {
		listener, err := net.Listen("tcp", ":9095")
		if err != nil {
			errors <- err
			return
		}
		grpcServer := grpc.NewServer()
		pb.RegisterPreviewImageServiceServer(grpcServer, transport.NewGRPCServer(ctx,
			endpoint.Endpoints{
				GetPreviewImageByUserEndpoint: endpoint.MakeGetPreviewImageByUserEndpoint(svc),
				SavePreviewImageEndpoint:      endpoint.MakeSavePreviewImageEndpoint(svc),
			}))
		logger.Log("msg", "PreviewImage gRPC server", "port", ":9095")
		errors <- grpcServer.Serve(listener)
	}()
	fmt.Println(<-errors)
}
