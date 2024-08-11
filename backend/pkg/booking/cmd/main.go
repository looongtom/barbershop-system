package main

import (
	"DoAn/database"
	"DoAn/pkg/booking"
	repository "DoAn/pkg/booking/db"
	"DoAn/pkg/booking/service"
	"DoAn/pkg/booking/transport"
	"google.golang.org/grpc"

	"fmt"
	logV "log"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}
	logger := log.NewLogfmtLogger(os.Stderr)
	//collectionMongo := database.ConnectMongo(os.Getenv("TokenCollectionMongo"))
	collectionPostgres, err := database.ConnectPostgresBooking()
	if err != nil {
		logV.Fatalf("Error getting env, %v", err)
	}
	err = booking.CreateTable(*collectionPostgres)
	if err != nil {
		logV.Fatalf("Error creating table: %v", err)
	}

	r := mux.NewRouter()

	var svc booking.BookingService
	svc = service.BookingStruct{}
	{
		repo, err := repository.NewRepository(collectionPostgres, logger)

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

		connGrpcTimeslit, err := grpc.Dial(os.Getenv("GRPC_TIMESLOT_SERVER"), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			fmt.Printf("did not connect: %v", err)
			logV.Fatalf("Error getting env, %v", err)
		}
		defer connGrpcTimeslit.Close()
		if err != nil {
			fmt.Printf("Error getting env, %v", err)
			logV.Fatalf("Error getting env, %v", err)
		}

		svc = service.NewService(repo, logger, connGrpcAccount, connGrpcTimeslit)
	}

	CreateBookingHandler := httptransport.NewServer(
		transport.MakeCreateBookingEndpoints(svc),
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
		transport.DecodeEmptyRequest,
		transport.EncodeResponse,
	)

	UpdateBookingHandler := httptransport.NewServer(
		transport.MakeUpdateBookingEndpoints(svc),
		transport.DecodeUpdateBookingRequest,
		transport.EncodeResponse,
	)

	http.Handle("/", addCorsHeaders(r))

	r.Handle("/booking/create", CreateBookingHandler).Methods("POST")
	r.Handle("/booking/update", UpdateBookingHandler).Methods("POST")
	r.Handle("/booking/get-by-id", GetBookingHandler).Methods("GET")
	r.Handle("/booking/get-list", GetListBookingHandler).Methods("GET")

	logger.Log("msg", "HTTP", "addr", ":8002")
	logger.Log("err", http.ListenAndServe(":8002", nil))
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
