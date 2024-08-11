package main

import (
	"DoAn/database"
	"DoAn/pkg/timeslot"
	repository "DoAn/pkg/timeslot/db"
	"DoAn/pkg/timeslot/endpoint"
	"DoAn/pkg/timeslot/pb"
	"DoAn/pkg/timeslot/service"
	"DoAn/pkg/timeslot/transport"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	logV "log"
	"net"
	"os"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}
	logger := log.NewLogfmtLogger(os.Stderr)
	collectionPostgres, err := database.ConnectPostgresTimeSlot()
	if err != nil {
		logV.Fatalf("Error getting env, %v", err)
	}
	ctx := context.Background()

	var svc timeslot.TimeSlotService
	svc = service.TimeSlotService{}
	{
		repo, err := repository.NewRepository(collectionPostgres, logger)
		if err != nil {
			logV.Fatalf("Error getting env, %v", err)
		}
		svc = service.NewService(repo, logger)
	}
	errors := make(chan error)
	go func() {
		listener, err := net.Listen("tcp", ":9093")
		if err != nil {
			errors <- err
			return
		}
		grpcServer := grpc.NewServer()
		pb.RegisterTimeslotServiceServer(grpcServer, transport.NewGRPCServer(ctx,
			endpoint.Endpoints{
				FindExistTimeslotEndpoint:    endpoint.MakeFindExistTimeslotEndpoint(svc),
				CheckExistTimeslotEndpoint:   endpoint.MakeCheckExistTimeslotEndpoint(svc),
				UpdateTimeslotStatusEndpoint: endpoint.MakeUpdateTimeslotStatusEndpoint(svc),
			}))
		logger.Log("msg", "Timeslot gRPC server", "port", ":9093")
		errors <- grpcServer.Serve(listener)
	}()
	fmt.Println(<-errors)
}
