package main

import (
	booking "DoAn"
	"DoAn/database"
	repository "DoAn/db"
	"DoAn/middleware"
	"DoAn/service"
	"DoAn/transport"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"fmt"
	logV "log"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	kafkaBrokerServer = "0.tcp.ap.ngrok.io:16436"
)

func main() {
	err := godotenv.Load("booking.env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}
	logger := log.NewLogfmtLogger(os.Stderr)
	kafkaBrokerServer = os.Getenv("KAFKA_BROKER")
	// collectionMongo := database.ConnectMongo(os.Getenv("TokenCollectionMongo"))
	collectionPostgres, err := database.ConnectPostgresBooking()
	if err != nil {
		logV.Fatalf("Error getting env, %v", err)
	}
	err = booking.CreateTable(*collectionPostgres)
	if err != nil {
		logV.Fatalf("Error creating table: %v", err)
	}
	r := mux.NewRouter()

	connGrpcAccount, err := grpc.NewClient(os.Getenv("GRPC_ACCOUNT_SERVER"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("did not connect: %v", err)
		logV.Fatalf("Error getting env, %v", err)
	}
	defer connGrpcAccount.Close()

	var svc booking.BookingService
	svc = service.BookingStruct{}
	{
		repo, err := repository.NewRepository(collectionPostgres, logger)
		if err != nil {
			logV.Fatalf("Error loading repository, %v", err)
		}

		connGrpcTimeslot, err := grpc.NewClient(os.Getenv("GRPC_TIMESLOT_SERVER"), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("did not connect: %v", err)
			logV.Fatalf("Error getting env, %v", err)
		}
		defer connGrpcTimeslot.Close()

		connGrpcService, err := grpc.NewClient(os.Getenv("GRPC_SERVICE_SERVER"), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("did not connect: %v", err)
			logV.Fatalf("Error getting env, %v", err)
		}
		defer connGrpcService.Close()

		kafkaBroker, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaBrokerServer})
		if err != nil {
			fmt.Printf("Failed to create producer: %s\n", err)
			return
		}
		defer kafkaBroker.Close()
		svc = service.NewService(repo, logger, connGrpcAccount, connGrpcTimeslot, connGrpcService, kafkaBroker)
	}

	CreateBookingHandler := httptransport.NewServer(
		transport.MakeCreateBookingEndpoints(svc),
		transport.DecodeCreateBookingRequest,
		transport.EncodeResponse,
	)

	CreateBookingKafkaHandler := httptransport.NewServer(
		transport.MakeCreateBookingKafkaEndpoints(svc),
		transport.DecodeCreateBookingRequest,
		transport.EncodeResponse,
	)

	GetBookingHandler := httptransport.NewServer(
		transport.MakeGetBookingEndpoints(svc),
		transport.DecodeGetBookingRequest,
		transport.EncodeResponse,
	)

	GetListBookingHandler := httptransport.NewServer(
		transport.MakeGetListBookingEndpoints(svc),
		transport.DecodeGetListBookingRequest,
		transport.EncodeResponse,
	)

	FindBookingHandler := httptransport.NewServer(
		transport.MakeFindBookingEndpoints(svc),
		transport.DecodeFindBookingRequest,
		transport.EncodeResponse,
	)

	UpdateBookingHandler := httptransport.NewServer(
		transport.MakeUpdateBookingEndpoints(svc),
		transport.DecodeUpdateBookingRequest,
		transport.EncodeResponse,
	)

	UpdateBookingServiceHandler := httptransport.NewServer(
		transport.MakeUpdateBookingServiceEndpoints(svc),
		transport.DecodeUpdateBookingServiceRequest,
		transport.EncodeResponse,
	)

	UpdateBookingTimeslotHandler := httptransport.NewServer(
		transport.MakeUpdateBookingTimeslotEndpoints(svc),
		transport.DecodeUpdateBookingTimeslotRequest,
		transport.EncodeResponse,
	)

	UpdateBookingStatusHandler := httptransport.NewServer(
		transport.MakeUpdateBookingStatusEndpoints(svc),
		transport.DecodeUpdateBookingStatusRequest,
		transport.EncodeResponse,
	)

	http.Handle("/", addCorsHeaders(r))

	r.Handle("/booking/create", middleware.JWTMiddleware(CreateBookingHandler, connGrpcAccount)).Methods("POST")
	r.Handle("/booking/create-kafka", middleware.JWTMiddleware(CreateBookingKafkaHandler, connGrpcAccount)).Methods("POST")
	r.Handle("/booking/update", middleware.JWTMiddlewareBarber(UpdateBookingHandler, connGrpcAccount)).Methods("POST")
	r.Handle("/booking/update-booking-service", middleware.JWTMiddlewareBarber(UpdateBookingServiceHandler, connGrpcAccount)).Methods("POST")
	r.Handle("/booking/update-booking-timeslot", middleware.JWTMiddlewareBarber(UpdateBookingTimeslotHandler, connGrpcAccount)).Methods("POST")
	r.Handle("/booking/update-booking-status", middleware.JWTMiddlewareBarber(UpdateBookingStatusHandler, connGrpcAccount)).Methods("POST")
	r.Handle("/booking/get-by-id", middleware.JWTMiddleware(GetBookingHandler, connGrpcAccount)).Methods("GET")
	r.Handle("/booking/get-list", middleware.JWTMiddleware(GetListBookingHandler, connGrpcAccount)).Methods("GET")
	r.Handle("/booking/find", middleware.JWTMiddlewareGetListBooking(FindBookingHandler, connGrpcAccount)).Methods("POST")

	logger.Log("msg", "HTTP", "addr", ":8010")
	logger.Log("err", http.ListenAndServe(":8010", nil))
}

func addCorsHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cho phép tất cả các origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Các headers được phép
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// Phương thức được phép
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,DELETE")

		// Nếu phương thức là OPTIONS, không cần xử lý tiếp
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Chuyển tiếp yêu cầu đến handler tiếp theo
		handler.ServeHTTP(w, r)
	})
}
