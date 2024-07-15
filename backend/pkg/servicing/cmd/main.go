package main

import (
	"DoAn/database"
	repository "DoAn/pkg/servicing/db"
	"DoAn/pkg/servicing/service"
	transport "DoAn/pkg/servicing/transport"
	"net/http"

	"DoAn/pkg/servicing"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	logV "log"
	"os"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	err := godotenv.Load(".env")
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

	http.Handle("/", addCorsHeaders(r))

	r.Handle("/servicing/category/create", CreateCategoryHandler).Methods("POST")
	r.Handle("/servicing/service/create", CreateServiceHandler).Methods("POST")
	r.Handle("/servicing/service/update", UpdateServiceHandler).Methods("POST")
	r.Handle("/servicing/category/update", UpdateCategoryHandler).Methods("POST")

	r.Handle("/servicing/category/get-list", GetListCategoryHandler).Methods("GET")
	r.Handle("/servicing/service/get-list", GetListServiceHandler).Methods("GET")
	r.Handle("/servicing/category/get", GetCategoryHandler).Methods("GET")
	r.Handle("/servicing/service/get", GetServiceHandler).Methods("GET")

	logger.Log("msg", "HTTP", "addr", ":8001")
	logger.Log("err", http.ListenAndServe(":8001", nil))

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
