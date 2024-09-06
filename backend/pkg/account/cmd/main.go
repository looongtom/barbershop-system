package main

import (
	"DoAn/pkg/account"
	"DoAn/pkg/account/auth"
	repository2 "DoAn/pkg/account/auth/repository"
	service2 "DoAn/pkg/account/auth/service"
	"DoAn/pkg/account/database"
	repository "DoAn/pkg/account/db"
	"DoAn/pkg/account/middleware"
	"DoAn/pkg/account/service"
	"DoAn/pkg/account/transport"
	"github.com/joho/godotenv"
	logV "log"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func main() {
	//read the account.env file
	err := godotenv.Load("account.env")
	if err != nil {
		logV.Fatalf("Error loading account.env file: %v", err)
	}
	logger := log.NewLogfmtLogger(os.Stderr)
	collectionMongo := database.ConnectMongo(os.Getenv("TokenCollectionMongo"))
	collectionMongoBlkToken := database.ConnectMongo(os.Getenv("TokenBlackListCollectionMongo"))
	collectionPostgres, err := database.ConnectPostgres()
	collectionRedis := auth.ConnectRedis(os.Getenv("REDIS_ADDRESS"), os.Getenv("REDIS_PASSWORD"))
	secretKey := []byte(os.Getenv("SECRET_JWT"))

	if err != nil {
		logV.Fatalf("Error getting env, %v", err)
	}
	r := mux.NewRouter()

	var authenSvc auth.AuthenService
	authenSvc = service2.AuthServiceStruct{}
	{
		repo2 := repository2.NewTokenRepository(collectionRedis)
		authenSvc = service2.NewAuthService(repo2, secretKey)
	}

	var svc account.UserService
	svc = service.UserServiceStruct{}
	{
		repo, err := repository.NewRepository(collectionMongo, collectionMongoBlkToken, collectionPostgres, logger)
		if err != nil {
			logV.Fatalf("Error getting env, %v", err)
		}
		repo2 := repository2.NewTokenRepository(collectionRedis)
		authenSvc = service2.NewAuthService(repo2, secretKey)
		svc = service.NewService(repo, authenSvc, logger)
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
	//ChangePassFirstTimeHandler := httptransport.NewServer(
	//	transport.MakeChangePassFirstTimeEndpoints(svc),
	//	transport.DecodeChangePassFirstTimeRequest,
	//	transport.EncodeResponse)
	LogoutHandler := httptransport.NewServer(
		transport.MakeLogoutEndpoints(authenSvc),
		transport.DecodeEmptyRequest,
		transport.EncodeResponse)
	RefreshHandler := httptransport.NewServer(
		transport.MakeRefreshEndpoints(authenSvc),
		transport.DecodeEmptyRequest,
		transport.EncodeResponse)

	http.Handle("/", addCorsHeaders(r))

	r.Handle("/auth/register", RegisterUserHandler).Methods("POST")
	r.Handle("/auth/login", LoginHandler).Methods("POST")
	r.Handle("/auth/refresh", middleware.JWTMiddlewareRefreshToken(RefreshHandler, authenSvc)).Methods("GET")
	r.Handle("/auth/profile", middleware.JWTMiddleware(GetProfileHandler, authenSvc)).Methods("GET")
	r.Handle("/auth/logout", middleware.JWTMiddleware(LogoutHandler, authenSvc)).Methods("POST")
	//r.Handle("/auth/change-pass-first-time", ChangePassFirstTimeHandler).Methods("POST")
	//
	//r.Handle("/login-oauth", http.HandlerFunc(transport.HandleMain)).Methods("GET")
	//r.Handle("/login", http.HandlerFunc(transport.HandleGoogleLogin)).Methods("GET")
	//r.Handle("/callback", http.HandlerFunc(transport.HandleGoogleCallback)).Methods("GET")

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
