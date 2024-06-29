package main

import (
	"DoAn/database"
	"DoAn/pkg/user"
	repository "DoAn/pkg/user/mongodb"
	"DoAn/pkg/user/service"
	"DoAn/pkg/user/transport"
	"github.com/joho/godotenv"
	logV "log"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}
	logger := log.NewLogfmtLogger(os.Stderr)
	collectionMongo := database.ConnectMongo(os.Getenv("TokenCollectionMongo"))
	collectionPostgres, err := database.ConnectPostgres()
	if err != nil {
		logV.Fatalf("Error getting env, %v", err)
	}
	r := mux.NewRouter()

	var svc user.UserService
	svc = service.UserServiceStruct{}
	{
		repo, err := repository.NewRepository(collectionMongo, collectionPostgres, logger)
		if err != nil {
			logV.Fatalf("Error getting env, %v", err)
		}
		svc = service.NewService(repo, logger)
	}

	RegisterUserHandler := httptransport.NewServer(
		transport.MakeRegisterUserEndpoints(svc),
		transport.DecodeRegisterUserRequest,
		transport.EncodeResponse)
	LoginHandler := httptransport.NewServer(
		transport.MakeLoginEndpoints(svc),
		transport.DecodeLoginRequest,
		transport.EncodeResponse)
	GetProfileHandler := httptransport.NewServer(
		transport.MakeGetProfileEndpoints(svc),
		transport.DecodeGetProfileRequest,
		transport.EncodeResponse)
	ChangePassFirstTimeHandler := httptransport.NewServer(
		transport.MakeChangePassFirstTimeEndpoints(svc),
		transport.DecodeChangePassFirstTimeRequest,
		transport.EncodeResponse)
	LogoutHandler := httptransport.NewServer(
		transport.MakeLogoutEndpoints(svc),
		transport.DecodeLogoutRequest,
		transport.EncodeResponse)
	RefreshHandler := httptransport.NewServer(
		transport.MakeRefreshEndpoints(svc),
		transport.DecodeRefreshRequest,
		transport.EncodeResponse)

	http.Handle("/", addCorsHeaders(r))

	r.Handle("/auth/register", RegisterUserHandler).Methods("POST")
	r.Handle("/auth/login", LoginHandler).Methods("POST")
	r.Handle("/auth/refresh", RefreshHandler).Methods("POST")
	r.Handle("/auth/profile", GetProfileHandler).Methods("GET")
	r.Handle("/auth/logout", LogoutHandler).Methods("POST")
	r.Handle("/auth/change-pass-first-time", ChangePassFirstTimeHandler).Methods("POST")

	r.Handle("/login-oauth", http.HandlerFunc(transport.HandleMain)).Methods("GET")
	r.Handle("/login", http.HandlerFunc(transport.HandleGoogleLogin)).Methods("GET")
	r.Handle("/callback", http.HandlerFunc(transport.HandleGoogleCallback)).Methods("GET")

	logger.Log("msg", "HTTP", "addr", ":8000")
	logger.Log("err", http.ListenAndServe(":8000", nil))
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
