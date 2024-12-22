package main

import (
	account "DoAn"
	"DoAn/auth"
	repository2 "DoAn/auth/repository"
	service2 "DoAn/auth/service"
	"DoAn/database"
	repository "DoAn/db"
	"DoAn/middleware"
	"DoAn/service"
	"DoAn/transport"

	"context"
	"fmt"
	logV "log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func scheduledTask(svc auth.AuthenService) {
	err := svc.CleanToken(context.Background())
	if err != nil {
		fmt.Printf("error while cleaning token: %v", err)
	}
	fmt.Println("Task executed at:", time.Now())
}

func scheduleDailyAt(hour, min, sec int, svc auth.AuthenService, task func(svc auth.AuthenService)) {
	for {
		now := time.Now()
		scheduledTime := time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, now.Location())
		if now.After(scheduledTime) {
			scheduledTime = scheduledTime.Add(24 * time.Hour)
		}
		duration := time.Until(scheduledTime)

		timer := time.NewTimer(duration)
		<-timer.C
		task(svc)
	}
}

func main() {
	// read the account.env file
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
	secretRefresh := []byte(os.Getenv("SECRET_REFRESH"))

	var scheduledTime time.Time
	scheduledTimeStr := os.Getenv("SCHEDULED_TIME")
	if scheduledTimeStr == "" {
		fmt.Println("SCHEDULED_TIME not set in .env file")
		scheduledTime = time.Now()
	}
	scheduledTime, err = time.Parse("15:04:05", scheduledTimeStr)
	if err != nil {
		fmt.Printf("Error parsing scheduled time: %v\n", err)
		return
	}
	// scheduledTime = time.Now().Add(5 * time.Second)
	fmt.Println("Scheduled time: ", scheduledTime)

	if err != nil {
		logV.Fatalf("Error getting env, %v", err)
	}

	err = account.CreateTable(*collectionPostgres)
	if err != nil {
		logV.Fatalf("Error creating table: %v", err)
	}

	r := mux.NewRouter()

	var authenSvc auth.AuthenService
	authenSvc = service2.AuthServiceStruct{}
	{
		repo2 := repository2.NewTokenRepository(collectionRedis)
		authenSvc = service2.NewAuthService(repo2, secretKey, secretRefresh)
		go scheduleDailyAt(scheduledTime.Hour(), scheduledTime.Minute(), scheduledTime.Second(), authenSvc, scheduledTask)
	}

	var svc account.UserService
	svc = service.UserServiceStruct{}
	{
		repo, err := repository.NewRepository(collectionMongo, collectionMongoBlkToken, collectionPostgres, logger)
		if err != nil {
			logV.Fatalf("Error getting env, %v", err)
		}
		repo2 := repository2.NewTokenRepository(collectionRedis)
		authenSvc = service2.NewAuthService(repo2, secretKey, secretRefresh)
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
	ChangePassFirstTimeHandler := httptransport.NewServer(
		transport.MakeChangePassFirstTimeEndpoints(svc),
		transport.DecodeChangePassFirstTimeRequest,
		transport.EncodeResponse)
	LogoutHandler := httptransport.NewServer(
		transport.MakeLogoutEndpoints(authenSvc),
		transport.DecodeEmptyRequest,
		transport.EncodeResponse)
	RefreshHandler := httptransport.NewServer(
		transport.MakeRefreshEndpoints(authenSvc),
		transport.DecodeEmptyRequest,
		transport.EncodeResponse)

	GetListBarberHandler := httptransport.NewServer(
		transport.MakeGetListBarberEndpoints(svc),
		transport.DecodeEmpty,
		transport.EncodeResponse)

	GoogleCallbackHandler := httptransport.NewServer(
		transport.MakeGoogleCallbackEndpoints(svc),
		transport.DecodeGoogleCallbackRequest,
		transport.EncodeResponse)

	http.Handle("/", addCorsHeaders(r))

	r.Handle("/auth/register", RegisterUserHandler).Methods("POST")
	r.Handle("/auth/login", LoginHandler).Methods("POST")

	r.Handle("/barber/get-list", GetListBarberHandler).Methods("GET")

	r.Handle("/auth/refresh", middleware.JWTMiddlewareRefreshToken(RefreshHandler, authenSvc)).Methods("GET")
	r.Handle("/auth/profile", middleware.JWTMiddleware(GetProfileHandler, authenSvc)).Methods("GET")
	r.Handle("/auth/logout", middleware.JWTMiddleware(LogoutHandler, authenSvc)).Methods("POST")
	r.Handle("/auth/change-pass-first-time", middleware.JWTMiddleware(ChangePassFirstTimeHandler, authenSvc)).Methods("POST")
	//
	r.Handle("/login-oauth", http.HandlerFunc(transport.HandleMain)).Methods("GET")
	r.Handle("/login", http.HandlerFunc(transport.HandleGoogleLogin)).Methods("GET")
	// r.Handle("/callback", http.HandlerFunc(transport.HandleGoogleCallback)).Methods("GET")
	r.Handle("/callback", GoogleCallbackHandler).Methods("GET")

	logger.Log("msg", "HTTP", "addr", ":8008")
	logger.Log("err", http.ListenAndServe(":8008", nil))

	// Keep the main function running
	select {}
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
