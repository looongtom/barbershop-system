package main

import (
	"DoAn/pkg/criteria"
	"DoAn/pkg/criteria/database"
	repository "DoAn/pkg/criteria/db"
	"DoAn/pkg/criteria/service"
	"DoAn/pkg/criteria/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	logV "log"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
)

func main() {
	err := godotenv.Load("criteria.env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}
	logger := log.NewLogfmtLogger(os.Stderr)
	collectionPostgres, err := database.ConnectPostgresCriteria()
	if err != nil {
		logV.Fatalf("Error getting env, %v", err)
	}
	err = criteria.CreateTable(*collectionPostgres)
	if err != nil {
		logV.Fatalf("Error creating table: %v", err)
	}
	r := mux.NewRouter()

	var svc criteria.CriteriaService
	svc = service.CriteriaStruct{}
	{
		repo, err := repository.NewRepository(collectionPostgres, logger)
		if err != nil {
			logV.Fatalf("Error getting env, %v", err)
		}
		svc = service.NewService(repo, logger)
	}
	CreateOrUpdateCategoryHandler := httptransport.NewServer(
		transport.MakeCreateCategoryEndpoints(svc),
		transport.DecodeCreateOrUpdateCategoryRequest,
		transport.EncodeResponse,
	)

	CreateOrUpdateCriteriaHandler := httptransport.NewServer(
		transport.MakeCreateOrUpdateCriteriaEndpoints(svc),
		transport.DecodeCreateOrUpdateCriteriaRequest,
		transport.EncodeResponse,
	)

	GetCriteriaHandler := httptransport.NewServer(
		transport.MakeGetCriteriaEndpoints(svc),
		transport.DecodeGetRequest,
		transport.EncodeResponse,
	)

	GetCategoryHandler := httptransport.NewServer(
		transport.MakeGetCategoryEndpoints(svc),
		transport.DecodeGetRequest,
		transport.EncodeResponse,
	)

	GetListCategoryHandler := httptransport.NewServer(
		transport.MakeGetListCategoryEndpoints(svc),
		transport.DecodeEmptyRequest,
		transport.EncodeResponse,
	)

	GetListCriteriaHandler := httptransport.NewServer(
		transport.MakeGetListCriteriaEndpoints(svc),
		transport.DecodeEmptyRequest,
		transport.EncodeResponse,
	)

	DeleteCategoryHandler := httptransport.NewServer(
		transport.MakeDeleteCategoryEndpoints(svc),
		transport.DecodeGetRequest,
		transport.EncodeResponse,
	)

	DeleteCriteriaHandler := httptransport.NewServer(
		transport.MakeDeleteCriteriaEndpoints(svc),
		transport.DecodeGetRequest,
		transport.EncodeResponse,
	)

	http.Handle("/", addCorsHeaders(r))

	r.Handle("/criteria/category/create-or-update", CreateOrUpdateCategoryHandler).Methods("POST")
	r.Handle("/criteria/service/create-or-update", CreateOrUpdateCriteriaHandler).Methods("POST")

	r.Handle("/criteria/category/get-list", GetListCategoryHandler).Methods("GET")
	r.Handle("/criteria/service/get-list", GetListCriteriaHandler).Methods("GET")
	r.Handle("/criteria/category/get", GetCategoryHandler).Methods("GET")
	r.Handle("/criteria/service/get", GetCriteriaHandler).Methods("GET")
	r.Handle("/criteria/category/delete", DeleteCategoryHandler).Methods("DELETE")
	r.Handle("/criteria/service/delete", DeleteCriteriaHandler).Methods("DELETE")

	logger.Log("msg", "HTTP", "addr", ":8006")
	logger.Log("err", http.ListenAndServe(":8006", nil))
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
