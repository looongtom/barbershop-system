package timeslot

import (
	"database/sql"
	"fmt"
)

func CreateTable(db sql.DB) (err error) {

	tableQueries := []string{
		`CREATE TABLE IF NOT EXISTS timeslot(
    			id SERIAL PRIMARY KEY,
    			barber_id int NOT NULL,
    			start_time VARCHAR(255) NOT NULL,
    			booked_date VARCHAR(255) NOT NULL,
    			status VARCHAR(255) NOT NULL,
    			created_at BIGINT NOT NULL,
    			updated_at BIGINT NOT NULL);
`,
	}

	for _, query := range tableQueries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}

	fmt.Println("Successfully created BookingDetail table!")
	return nil
}
