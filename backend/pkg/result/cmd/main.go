package main

import (
	"DoAn/pkg/result"
	"DoAn/pkg/result/database"
	repository "DoAn/pkg/result/db"
	"DoAn/pkg/result/middleware"
	"DoAn/pkg/result/service"
	"DoAn/pkg/result/transport"
	"fmt"
	logV "log"
	"net/http"
	"os"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-kit/kit/log"
)

func main() {
	err := godotenv.Load("result.env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}
	logger := log.NewLogfmtLogger(os.Stderr)
	collectionPostgres, err := database.ConnectPostgresResult()
	if err != nil {
		logV.Fatalf("Error getting env, %v", err)
	}
	err = result.CreateTable(*collectionPostgres)
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

	var svc result.ResultService
	svc = service.ResultService{}
	{
		repo, err := repository.NewRepository(collectionPostgres, logger)
		if err != nil {
			logV.Fatalf("Error getting env, %v", err)
		}
		svc = service.NewService(repo, logger)
	}

	CreateOrUpdateResultHandler := httptransport.NewServer(
		transport.MakeCreateOrUpdateResultEndpoints(svc),
		transport.DecodeCreateOrUpdateResult,
		transport.EncodeResponse,
	)

	GetResultByBarberIdHandler := httptransport.NewServer(
		transport.MakeGetResultByBarberIdEndpoints(svc),
		transport.DecodeGetResultById,
		transport.EncodeResponse,
	)

	GetResultByUserIdHandler := httptransport.NewServer(
		transport.MakeGetResultByUserIdEndpoints(svc),
		transport.DecodeGetResultById,
		transport.EncodeResponse,
	)

	GetResultByBookingIdHandler := httptransport.NewServer(
		transport.MakeGetResultByBookingIdEndpoints(svc),
		transport.DecodeGetResultById,
		transport.EncodeResponse,
	)

	http.Handle("/", addCorsHeaders(r))

	r.Methods("POST").Path("/createOrUpdateResult").Handler(middleware.JWTMiddlewareBarber(CreateOrUpdateResultHandler, connGrpcAccount))
	r.Methods("GET").Path("/getResultByBarberId").Handler(middleware.JWTMiddleware(GetResultByBarberIdHandler, connGrpcAccount))
	r.Methods("GET").Path("/getResultByUserId").Handler(middleware.JWTMiddleware(GetResultByUserIdHandler, connGrpcAccount))
	r.Methods("GET").Path("/getResultByBookingId").Handler(middleware.JWTMiddleware(GetResultByBookingIdHandler, connGrpcAccount))

	logger.Log("msg", "HTTP", "addr", ":8007")
	logger.Log("err", http.ListenAndServe(":8007", nil))

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
