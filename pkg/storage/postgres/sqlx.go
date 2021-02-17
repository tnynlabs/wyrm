package postgres

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgresql driver
)

// GetFromEnv Get postgresql db connection from enviornment variables
// Example:
// 		DB_HOST=localhost
// 		DB_USER=admin
// 		DB_PASSWORD=admin
// 		DB_NAME=dev
// 		DB_PORT=5432
// 		DB_SCHEMA_PATH=cmd/db-init/schema.sql
func GetFromEnv() (*sqlx.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}
