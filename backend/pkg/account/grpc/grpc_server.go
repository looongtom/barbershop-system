package main

import (
	"DoAn"
	"DoAn/auth"
	repository2 "DoAn/auth/repository"
	service2 "DoAn/auth/service"
	"DoAn/database"
	repository "DoAn/db"
	"DoAn/endpoint"
	"DoAn/pb"
	"DoAn/service"
	"DoAn/transport"
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
	collectionRedis := auth.ConnectRedis(os.Getenv("REDIS_ADDRESS"), os.Getenv("REDIS_PASSWORD"))
	secretKey := []byte(os.Getenv("SECRET_JWT"))
	secretRefresh := []byte(os.Getenv("SECRET_REFRESH"))
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

	var authenSvc auth.AuthenService
	authenSvc = service2.AuthServiceStruct{}
	{
		repo2 := repository2.NewTokenRepository(collectionRedis)
		authenSvc = service2.NewAuthService(repo2, secretKey, secretRefresh)
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
			endpoint.Endpoints{
				CheckExistedBarberEndpoint: endpoint.MakeCheckExistedBarberEndpoint(svc),
				VerifyTokenEndpoint:        endpoint.MakeVerifyTokenEndpoint(authenSvc),
			}))

		logger.Log("msg", "Account gRPC server", "port", ":9090")
		errors <- gRPCServer.Serve(listener)
	}()

	fmt.Println(<-errors)
}
