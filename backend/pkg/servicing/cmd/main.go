package main

import (
	"DoAn/database"
	repository "DoAn/db"
	"DoAn/middleware"
	"DoAn/service"
	transport "DoAn/transport"
	"fmt"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"DoAn"
	logV "log"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	err := godotenv.Load("servicing.env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}
	logger := log.NewLogfmtLogger(os.Stderr)
	//collectionMongo := database.ConnectMongo(os.Getenv("TokenCollectionMongo"))
	collectionPostgres, err := database.ConnectPostgresServicing()
	if err != nil {
		logV.Fatalf("Error getting env, %v", err)
	}
	r := mux.NewRouter()

	connGrpcAccount, err := grpc.NewClient(os.Getenv("GRPC_ACCOUNT_SERVER"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("did not connect: %v", err)
		logV.Fatalf("Error getting env, %v", err)
	}
	defer connGrpcAccount.Close()

	var svc servicing.ServicingService
	svc = service.ServicingStruct{}
	{
		repo, err := repository.NewRepository(collectionPostgres, logger)
		if err != nil {
			logV.Fatalf("Error getting env, %v", err)
		}
		svc = service.NewService(repo, logger)
	}

	CreateCategoryHandler := httptransport.NewServer(
		transport.MakeCreateCategoryEndpoints(svc),
		transport.DecodeCreateCategoryRequest,
		transport.EncodeResponse,
	)

	CreateServiceHandler := httptransport.NewServer(
		transport.MakeCreateServiceEndpoints(svc),
		transport.DecodeCreateServiceRequest,
		transport.EncodeResponse,
	)

	GetCategoryHandler := httptransport.NewServer(
		transport.MakeGetServiceEndpoints(svc),
		transport.DecodeGetServiceRequest,
		transport.EncodeResponse,
	)

	GetServiceHandler := httptransport.NewServer(
		transport.MakeGetServiceEndpoints(svc),
		transport.DecodeGetServiceRequest,
		transport.EncodeResponse,
	)

	GetListCategoryHandler := httptransport.NewServer(
		transport.MakeGetListCategoryEndpoints(svc),
		transport.DecodeEmptyRequest,
		transport.EncodeResponse,
	)

	GetListServiceHandler := httptransport.NewServer(
		transport.MakeGetListServiceEndpoints(svc),
		transport.DecodeEmptyRequest,
		transport.EncodeResponse,
	)

	UpdateServiceHandler := httptransport.NewServer(
		transport.MakeUpdateServiceEndpoints(svc),
		transport.DecodeCreateServiceRequest,
		transport.EncodeResponse,
	)

	UpdateCategoryHandler := httptransport.NewServer(
		transport.MakeUpdateCategoryEndpoints(svc),
		transport.DecodeUpdateCategoryRequest,
		transport.EncodeResponse,
	)

	GetListServiceAndCategoryHandler := httptransport.NewServer(
		transport.MakeGetListServiceAndCategoryEndpoints(svc),
		transport.DecodeEmptyRequest,
		transport.EncodeResponse,
	)

	http.Handle("/", addCorsHeaders(r))

	r.Handle("/servicing/category/create", middleware.JWTMiddlewareAdmin(CreateCategoryHandler, connGrpcAccount)).Methods("POST")
	r.Handle("/servicing/service/create", middleware.JWTMiddlewareAdmin(CreateServiceHandler, connGrpcAccount)).Methods("POST")
	r.Handle("/servicing/service/update", middleware.JWTMiddlewareAdmin(UpdateServiceHandler, connGrpcAccount)).Methods("POST")
	r.Handle("/servicing/category/update", middleware.JWTMiddlewareAdmin(UpdateCategoryHandler, connGrpcAccount)).Methods("POST")

	r.Handle("/servicing/category/get-list", middleware.JWTMiddleware(GetListCategoryHandler, connGrpcAccount)).Methods("GET")
	r.Handle("/servicing/service/get-list", middleware.JWTMiddleware(GetListServiceHandler, connGrpcAccount)).Methods("GET")
	r.Handle("/servicing/service/get-list-v2", middleware.JWTMiddleware(GetListServiceAndCategoryHandler, connGrpcAccount)).Methods("GET")
	r.Handle("/servicing/service/get-list-test", GetListServiceAndCategoryHandler).Methods("GET")

	r.Handle("/servicing/category/get", middleware.JWTMiddleware(GetCategoryHandler, connGrpcAccount)).Methods("GET")
	r.Handle("/servicing/service/get", middleware.JWTMiddleware(GetServiceHandler, connGrpcAccount)).Methods("GET")

	logger.Log("msg", "HTTP", "addr", ":8009")
	logger.Log("err", http.ListenAndServe(":8009", nil))

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
