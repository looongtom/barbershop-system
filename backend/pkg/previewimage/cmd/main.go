package main

import (
	"DoAn/database"
	"DoAn/pkg/previewimage"
	repository "DoAn/pkg/previewimage/db"
	"DoAn/pkg/previewimage/service"
	"DoAn/pkg/previewimage/transport"
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
	collectionPostgres, err := database.ConnectPostgresPreviewImage()
	if err != nil {
		logV.Fatalf("Error getting env, %v", err)
	}
	err = previewimage.CreateTablePreviewImage(*collectionPostgres)
	if err != nil {
		logV.Fatalf("Error creating table: %v", err)
	}
	r := mux.NewRouter()

	var svc previewimage.PreviewImageService
	svc = service.PreviewImageService{}
	{
		repo, err := repository.NewRepository(collectionPostgres, logger)
		if err != nil {
			logV.Fatalf("Error loading repository, %v", err)
		}
		svc = service.NewService(repo, logger)
	}

	CreatePreviewImageHandler := httptransport.NewServer(
		transport.MakeCreatePreviewImageEndpoints(svc),
		transport.DecodeCreatePreviewImageRequest,
		transport.EncodeResponse,
	)

	http.Handle("/", addCorsHeaders(r))

	r.Handle("/previewimage/create", CreatePreviewImageHandler).Methods("POST")

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
