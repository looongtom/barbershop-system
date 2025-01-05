package booking

import (
	"database/sql"
	"fmt"
)

func CreateTable(db sql.DB) (err error) {

	tableQueries := []string{
		`CREATE TABLE IF NOT EXISTS booking(
    		id SERIAL PRIMARY KEY,
    		customer_id int NOT NULL,
    		barber_id int NOT NULL,
    		result_id int,
    		status VARCHAR(255) NOT NULL,
    		price float NOT NULL,
    		slot_id int NOT NULL,
    		feedback_id int ,
    		preview_id int,
    		created_at BIGINT NOT NULL,
    		updated_at BIGINT NOT NULL);
`,
		`CREATE TABLE IF NOT EXISTS booking_detail (
			booking_id INT NOT NULL,
			service_id INT NOT NULL,
			FOREIGN KEY (booking_id) REFERENCES Booking(ID) ON DELETE CASCADE
		);`,
	}

	for _, query := range tableQueries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}
	return nil
}

func AlterTableAddColumn(db sql.DB) error {
	query := `ALTER TABLE booking ADD COLUMN booked_date VARCHAR(255);`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to alter table: %v", err)
	}
	return nil
}
