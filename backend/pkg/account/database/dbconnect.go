package database

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"fmt"
	"os"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(collectionName string) *mongo.Collection {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Connected to MongoDB!")
	collection := client.Database(os.Getenv("DATABASE_MONGO")).Collection(collectionName)

	return collection
}

func ConnectPostgres() (*sql.DB, error) {
	portString := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(portString)
	if err != nil {
		fmt.Println("Error connect postgres:", err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		port,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME_ACCOUNT"))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to PostgresSQL!")

	return db, nil

}
