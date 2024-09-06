package main

import (
	"DoAn/pkg/account"
	"DoAn/pkg/account/database"
	repository "DoAn/pkg/account/db"
	"DoAn/pkg/account/endpoint"
	"DoAn/pkg/account/pb"
	"DoAn/pkg/account/service"
	"DoAn/pkg/account/transport"
	"context"
	"fmt"
	logV "log"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/go-kit/kit/log"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("account.env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}
	logger := log.NewLogfmtLogger(os.Stderr)
	collectionMongo := database.ConnectMongo(os.Getenv("TokenCollectionMongo"))
	collectionPostgres, err := database.ConnectPostgres()
	if err != nil {
		logV.Fatalf("Error getting env, %v", err)
	}
	ctx := context.Background()

	var svc account.UserService
	svc = service.UserServiceStruct{}
	{
		repo, err := repository.NewRepository(collectionMongo, nil, collectionPostgres, logger)
		if err != nil {
			logV.Fatalf("Error getting env, %v", err)
		}
		svc = service.NewService(repo, nil, logger)
	}
	errors := make(chan error)
	go func() {
		listener, err := net.Listen("tcp", ":9090")
		if err != nil {
			errors <- err
			return
		}
		gRPCServer := grpc.NewServer()
		pb.RegisterUserServiceServer(gRPCServer, transport.NewGRPCServer(ctx,
			endpoint.Endpoints{CheckExistedBarberEndpoint: endpoint.MakeCheckExistedBarberEndpoint(svc)}))

		logger.Log("msg", "Account gRPC server", "port", ":9090")
		errors <- gRPCServer.Serve(listener)
	}()

	fmt.Println(<-errors)
}
