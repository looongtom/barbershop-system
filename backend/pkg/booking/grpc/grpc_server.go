package main

import (
	"DoAn"
	"DoAn/database"
	"DoAn/endpoint"

	repository "DoAn/db"
	"DoAn/pb"
	"DoAn/service"
	"DoAn/transport"
	"context"
	"fmt"
	logV "log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/go-kit/kit/log"
)

func main() {
	err := godotenv.Load("account.env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}
	logger := log.NewLogfmtLogger(os.Stderr)
	collectionPostgres, err := database.ConnectPostgresBooking()
	if err != nil {
		logV.Fatalf("Error getting env, %v", err)
	}
	ctx := context.Background()

	var svc booking.BookingService
	svc = service.BookingStruct{}
	{
		repo, err := repository.NewRepository(collectionPostgres, logger)
		if err != nil {
			logV.Fatalf("Error getting env, %v", err)
		}
		connGrpcAccount, err := grpc.Dial(os.Getenv("GRPC_ACCOUNT_SERVER"), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			fmt.Printf("did not connect: %v", err)
			logV.Fatalf("Error getting env, %v", err)
		}
		defer connGrpcAccount.Close()
		if err != nil {
			fmt.Printf("Error getting env, %v", err)
			logV.Fatalf("Error getting env, %v", err)
		}

		connGrpcTimeslot, err := grpc.Dial(os.Getenv("GRPC_TIMESLOT_SERVER"), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			fmt.Printf("did not connect: %v", err)
			logV.Fatalf("Error getting env, %v", err)
		}
		defer connGrpcTimeslot.Close()
		if err != nil {
			fmt.Printf("Error getting env, %v", err)
			logV.Fatalf("Error getting env, %v", err)
		}

		connGrpcService, err := grpc.Dial(os.Getenv("GRPC_SERVICE_SERVER"), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			fmt.Printf("did not connect: %v", err)
			logV.Fatalf("Error getting env, %v", err)
		}
		defer connGrpcService.Close()
		if err != nil {
			fmt.Printf("Error getting env, %v", err)
			logV.Fatalf("Error getting env, %v", err)
		}
		svc = service.NewService(repo, logger, connGrpcAccount, connGrpcTimeslot, connGrpcService, nil)
	}
	errors := make(chan error)
	go func() {
		listener, err := net.Listen("tcp", ":9094")
		if err != nil {
			errors <- err
			return
		}
		grpcServer := grpc.NewServer()
		pb.RegisterBookingServiceServer(grpcServer, transport.NewGRPCServer(ctx,
			endpoint.Endpoints{CreateBooking: endpoint.MakeCreateBookingEndpoint(svc)}))
		logger.Log("msg", "Booking gRPC server", "port", ":9094")
		errors <- grpcServer.Serve(listener)
	}()
	fmt.Println(<-errors)
}
