package main

import (
	"DoAn"
	"DoAn/database"
	repository "DoAn/db"
	"DoAn/endpoint"
	"DoAn/pb"
	"DoAn/service"
	"DoAn/transport"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"context"
	logV "log"
	"os"

	"github.com/joho/godotenv"

	"github.com/go-kit/kit/log"
)

func main() {
	err := godotenv.Load("account.env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}
	logger := log.NewLogfmtLogger(os.Stderr)
	collectionPostgres, err := database.ConnectPostgresServicing()
	if err != nil {
		logV.Fatalf("Error getting env, %v", err)
	}
	ctx := context.Background()
	var svc servicing.ServicingService
	svc = service.ServicingStruct{}
	{
		repo, err := repository.NewRepository(collectionPostgres, logger)
		if err != nil {
			logV.Fatalf("Error getting env, %v", err)
		}
		svc = service.NewService(repo, logger)
	}
	errors := make(chan error)
	go func() {
		listener, err := net.Listen("tcp", ":9091")
		if err != nil {
			errors <- err
			return
		}
		grpcServer := grpc.NewServer()
		pb.RegisterServicingServiceServer(grpcServer, transport.NewGRPCServer(ctx,
			endpoint.Endpoints{
				GetServiceById: endpoint.MakeGetServiceByIdEndpoint(svc),
			}))
		logger.Log("msg", "Service gRPC server", "port", ":9091")

		//fmt.Println("Service gRPC server is running on port 9091")
		errors <- grpcServer.Serve(listener)
	}()
	fmt.Println(<-errors)
}
