package result

import (
	"database/sql"
	"fmt"
)

func CreateTable(db sql.DB) (err error) {

	tableQueries := []string{
		`CREATE TABLE IF NOT EXISTS result(
    			id SERIAL PRIMARY KEY,
    			barber_id int NOT NULL,
    			user_id int NOT NULL,
    			booking_id int NOT NULL,
    			description VARCHAR(255) NOT NULL,
    			description VARCHAR(255) NOT NULL,
    			description VARCHAR(255) NOT NULL,
    			description VARCHAR(255) NOT NULL,
    			description VARCHAR(255) NOT NULL,
    			created_at BIGINT NOT NULL,
    			updated_at BIGINT NOT NULL);
`,
		`CREATE TABLE IF NOT EXISTS image_result(
    			id SERIAL PRIMARY KEY,
    			url VARCHAR(255) NOT NULL,
    			result_id int NOT NULL,
    			created_at BIGINT NOT NULL,
    			FOREIGN KEY (result_id) REFERENCES result(id) ON DELETE CASCADE
    );
`,
	}

	for _, query := range tableQueries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}
	fmt.Println("Successfully created Result table!")
	return nil
}
