package database

import (
	"database/sql"
	"fmt"
	"os"

	"orders-service/utils"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	var dsn string
	useCloudSQL := os.Getenv("USE_CLOUD_SQL") == "true"
	secrets, e := utils.GetSecrets([]string{"DB_PASSWORD", "DB_NAME", "DB_USER", "CLOUDSQL_CONNECTION_NAME"}, os.Getenv("PUBSUB_PROJECT_ID"))
	if e != nil {
		return fmt.Errorf("failed to get secrets: %w", e)
	}
	if useCloudSQL {
		// Use Cloud SQL Unix socket connection
		dsn = fmt.Sprintf("user=%s password=%s dbname=%s host=/cloudsql/%s sslmode=disable",
			secrets["DB_USER"],
			secrets["DB_PASSWORD"],
			secrets["DB_NAME"],
			secrets["CLOUDSQL_CONNECTION_NAME"],
		)
	} else {
		// Use local TCP connection
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			secrets["DB_USER"],
			secrets["DB_PASSWORD"],
			secrets["DB_NAME"])
	}

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open DB connection: %w", err)
	}

	return nil
}
