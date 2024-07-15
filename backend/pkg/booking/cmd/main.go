package cmd

import (
	"DoAn/database"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	logV "log"
	"os"

	"github.com/go-kit/kit/log"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		logV.Fatalln("Error getting env, %v", err)
	}
	logger := log.NewLogfmtLogger(os.Stderr)
	//collectionMongo := database.ConnectMongo(os.Getenv("TokenCollectionMongo"))
	collectionPostgres, err := database.ConnectPostgresBooking()
	if err != nil {
		logV.Fatalf("Error getting env, %v", err)
	}
	r := mux.NewRouter()


}
