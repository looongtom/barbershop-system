package main

import (
	"DoAn/database"
	"DoAn/pkg/criteria"
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
	err := godotenv.Load(".env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}
	logger := log.NewLogfmtLogger(os.Stderr)
	collectionPostgres, err := database.ConnectPostgresCriteria()
	if err != nil {
		logV.Fatalf("Error getting env, %v", err)
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
	CreateCategoryHandler := httptransport.NewServer(
		transport.MakeCreateCategoryEndpoints(svc),
		transport.DecodeCreateCategoryRequest,
		transport.EncodeResponse,
	)

	http.Handle("/", addCorsHeaders(r))

	r.Handle("/criteria/category/create", CreateCategoryHandler).Methods("POST")
	//r.Handle("/criteria/service/create", CreateCriteriaHandler).Methods("POST")
	//r.Handle("/criteria/service/update", UpdateCriteriaHandler).Methods("POST")
	//r.Handle("/criteria/category/update", UpdateCategoryHandler).Methods("POST")
	//
	//r.Handle("/criteria/category/get-list", GetListCategoryHandler).Methods("GET")
	//r.Handle("/criteria/service/get-list", GetListCriteriaHandler).Methods("GET")
	//r.Handle("/criteria/category/get", GetCategoryHandler).Methods("GET")
	//r.Handle("/criteria/service/get", GetCriteriaHandler).Methods("GET")

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
