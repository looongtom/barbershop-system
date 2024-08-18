package previewimage

import (
	"database/sql"
	"fmt"
)

func CreateTablePreviewImage(db sql.DB) (err error) {

	tableQueries := []string{
		`CREATE TABLE IF NOT EXISTS preview_image(
			id SERIAL PRIMARY KEY,
			account_id int NOT NULL,
			image_url VARCHAR(255) NOT NULL,
			created_at BIGINT NOT NULL);
`,
		`CREATE TABLE IF NOT EXISTS preview_image_criteria(
    		preview_image_id INT NOT NULL,
    		criteria_id INT NOT NULL,
    		FOREIGN KEY (preview_image_id) REFERENCES preview_image(ID) ON DELETE CASCADE
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
