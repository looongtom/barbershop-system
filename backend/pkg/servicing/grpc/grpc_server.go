package main

import (
	"DoAn/database"
	"DoAn/pkg/servicing"
	repository "DoAn/pkg/servicing/db"
	"DoAn/pkg/servicing/endpoint"
	"DoAn/pkg/servicing/pb"
	"DoAn/pkg/servicing/service"
	"DoAn/pkg/servicing/transport"
	"fmt"
	"google.golang.org/grpc"
	"net"

	"context"
	"github.com/joho/godotenv"
	logV "log"
	"os"

	"github.com/go-kit/kit/log"
)

func main() {
	err := godotenv.Load(".env")
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
		fmt.Println("gRPC server is running on port 9093")
		errors <- grpcServer.Serve(listener)
	}()
	fmt.Println(<-errors)
}
