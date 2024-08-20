package main

import (
	"DoAn/database"
	"DoAn/pkg/result"
	repository "DoAn/pkg/result/db"
	"DoAn/pkg/result/service"
	"DoAn/pkg/result/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	logV "log"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
)

func main() {
	err := godotenv.Load(".env")
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

	r.Methods("POST").Path("/createOrUpdateResult").Handler(CreateOrUpdateResultHandler)
	r.Methods("GET").Path("/getResultByBarberId").Handler(GetResultByBarberIdHandler)
	r.Methods("GET").Path("/getResultByUserId").Handler(GetResultByUserIdHandler)
	r.Methods("GET").Path("/getResultByBookingId").Handler(GetResultByBookingIdHandler)

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
