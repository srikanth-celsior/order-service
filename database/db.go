package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	var dsn string
	useCloudSQL := os.Getenv("USE_CLOUD_SQL") == "true"

	if useCloudSQL {
		// Use Cloud SQL Unix socket connection
		dsn = fmt.Sprintf("user=%s password=%s dbname=%s host=/cloudsql/%s sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("CLOUDSQL_CONNECTION_NAME"),
		)
	} else {
		// Use local TCP connection
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"))
	}

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open DB connection: %w", err)
	}

	return nil
}
