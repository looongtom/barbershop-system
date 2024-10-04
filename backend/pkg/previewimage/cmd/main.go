package main

import (
	"fmt"
	logV "log"
	"net/http"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"DoAn"
	"DoAn/database"
	repository "DoAn/db"
	"DoAn/middleware"
	"DoAn/service"
	"DoAn/transport"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-kit/kit/log"
)

const (
	kafkaBroker = "localhost:9092"
)

func main() {
	err := godotenv.Load("previewimage.env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}
	logger := log.NewLogfmtLogger(os.Stderr)
	collectionPostgres, err := database.ConnectPostgresPreviewImage()
	if err != nil {
		logV.Fatalf("Error getting env, %v", err)
	}
	err = previewimage.CreateTablePreviewImage(*collectionPostgres)
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

	kafkaBroker, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaBroker})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return
	}
	defer kafkaBroker.Close()

	var svc previewimage.PreviewImageService
	svc = service.PreviewImageService{}
	{
		repo, err := repository.NewRepository(collectionPostgres, logger)
		if err != nil {
			logV.Fatalf("Error loading repository, %v", err)
		}

		svc = service.NewService(repo, logger, connGrpcAccount, kafkaBroker)
	}

	CreatePreviewImageHandler := httptransport.NewServer(
		transport.MakeCreatePreviewImageEndpoints(svc),
		transport.DecodeCreatePreviewImageRequest,
		transport.EncodeResponse,
	)

	UploadImagesHandler := httptransport.NewServer(
		transport.MakeUploadImagesEndpoints(svc),
		transport.DecodeUploadImagesRequest,
		transport.EncodeResponse,
	)

	http.Handle("/", addCorsHeaders(r))

	r.Handle("/previewimage/create", middleware.JWTMiddleware(CreatePreviewImageHandler, connGrpcAccount)).Methods("POST")
	r.Handle("/previewimage/upload", middleware.JWTMiddleware(UploadImagesHandler, connGrpcAccount)).Methods("POST")

	logger.Log("msg", "HTTP", "addr", ":8005")
	logger.Log("err", http.ListenAndServe(":8005", nil))
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
