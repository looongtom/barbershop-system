package main

import (
	"fmt"
	logV "log"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	notification "DoAn"
	"DoAn/database"
	repository "DoAn/db"
	"DoAn/middleware"
	"DoAn/service"
	"DoAn/transport"
)

func main() {
	err := godotenv.Load("notification.env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}
	logger := log.NewLogfmtLogger(os.Stderr)
	collectionMongo := database.ConnectMongo(os.Getenv("NotificationCollectionMongo"))
	r := mux.NewRouter()

	connGrpcAccount, err := grpc.NewClient(os.Getenv("GRPC_ACCOUNT_SERVER"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("did not connect: %v", err)
		logV.Fatalf("Error getting env, %v", err)
	}
	defer connGrpcAccount.Close()

	var svc notification.NotificationService
	svc = service.NotificationStruct{}
	{
		repo := repository.NewRepository(collectionMongo, logger)
		svc = service.NewService(repo, logger)
	}

	GetListNotificationHandler := httptransport.NewServer(
		transport.MakeGetListNotificationEndpoints(svc),
		transport.DecodeGetListNotificationRequest,
		transport.EncodeResponse,
	)

	http.Handle("/", addCorsHeaders(r))

	r.Handle("/notification/get-list", middleware.JWTMiddleware(GetListNotificationHandler, connGrpcAccount)).Methods("POST")

	logger.Log("msg", "HTTP", "addr", ":8011")
	logger.Log("err", http.ListenAndServe(":8011", nil))

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
