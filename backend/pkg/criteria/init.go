package criteria

import (
	"database/sql"
	"fmt"
)

func CreateTable(db sql.DB) (err error) {

	tableQueries := []string{
		`CREATE TABLE IF NOT EXISTS criteria(
    			id SERIAL PRIMARY KEY,
    			name VARCHAR(255) NOT NULL,
    			img VARCHAR(255) NOT NULL,
    			category_id int NOT NULL,
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

	fmt.Println("Successfully created Criteria table!")
	return nil
}
