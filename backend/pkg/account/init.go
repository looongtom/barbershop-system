package account

import (
	"database/sql"
	"fmt"
)

func CreateTable(db sql.DB) (err error) {

	tableQueries := []string{
		`
		CREATE TABLE IF NOT EXISTS account (
			id SERIAL PRIMARY KEY,
            username VARCHAR(255) NOT NULL ,
            email VARCHAR(255) NOT NULL,
            password VARCHAR(255) NOT NULL,
		    role int NOT NULL ,
		    phone_number VARCHAR(255) NOT NULL,
			full_name VARCHAR(255) NOT NULL,
			about VARCHAR(255) NOT NULL,
			avatar VARCHAR(255) NOT NULL,
		    dob VARCHAR(255) NOT NULL,
			created_at BIGINT NOT NULL,
			updated_at BIGINT NOT NULL
		);
`,
	}
	for _, query := range tableQueries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}

	fmt.Println("Successfully created Account table!")
	return nil
}
