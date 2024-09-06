package database

import (
	"database/sql"
	_ "github.com/lib/pq"

	"fmt"
	"os"
	"strconv"
)

func ConnectPostgresBooking() (*sql.DB, error) {
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
		os.Getenv("DB_NAME_BOOKING"))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error connect postgres:", err)
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Error connect postgres:", err)
		panic(err)
	}

	fmt.Println("Successfully connected to PostgresSQL!")

	return db, nil

}
