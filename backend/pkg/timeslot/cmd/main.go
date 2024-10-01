package main

import (
	"DoAn/pkg/timeslot/database"
	"DoAn/pkg/timeslot/middleware"
	"DoAn/pkg/timeslot/transport"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"DoAn/pkg/timeslot"
	repository "DoAn/pkg/timeslot/db"
	"DoAn/pkg/timeslot/service"
	"fmt"
	logV "log"
	"net/http"
	"os"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/go-kit/kit/log"
)

func main() {
	err := godotenv.Load("timeslot.env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}
	logger := log.NewLogfmtLogger(os.Stderr)
	collectionPostgres, err := database.ConnectPostgresTimeSlot()
	if err != nil {
		logV.Fatalf("Error getting env, %v", err)
	}
	err = timeslot.CreateTable(*collectionPostgres)
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

	var svc timeslot.TimeSlotService
	svc = service.TimeSlotService{}
	{
		repo, err := repository.NewRepository(collectionPostgres, logger)
		if err != nil {
			fmt.Printf("Error getting env, %v", err)
			logV.Fatalf("Error getting env, %v", err)
		}
		svc = service.NewService(repo, logger)
	}

	CreateOrUpdateTimeslotHandler := httptransport.NewServer(
		transport.MakeCreateTimeSlotEndpoints(svc),
		transport.DecodeCreateTimeSlotRequest,
		transport.EncodeResponse,
	)

	CreateListTimeslotHandler := httptransport.NewServer(
		transport.MakeCreateListTimeSlotEndpoints(svc),
		transport.DecodeCreateListTimeSlotRequest,
		transport.EncodeResponse,
	)

	//UpdateTimeslotHandler := httptransport.NewServer(
	//	transport.MakeCreateTimeSlotEndpoints(svc),
	//	transport.DecodeUpdatedTimeSlotRequest,
	//	transport.EncodeResponse,
	//)

	GetListTimeSlotByBarberIdHandler := httptransport.NewServer(
		transport.MakeGetTimeSlotEndpoints(svc),
		transport.DecodeGetTimeSlotRequest,
		transport.EncodeResponse,
	)
	http.Handle("/", addCorsHeaders(r))

	r.Handle("/timeslot/find", GetListTimeSlotByBarberIdHandler).Methods("POST")
	r.Handle("/timeslot/create-or-update", middleware.JWTMiddlewareAdmin(CreateOrUpdateTimeslotHandler, connGrpcAccount)).Methods("POST")
	r.Handle("/timeslot/create-by-list", CreateListTimeslotHandler).Methods("POST")

	logger.Log("msg", "HTTP", "addr", ":8003")
	logger.Log("err", http.ListenAndServe(":8003", nil))
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
