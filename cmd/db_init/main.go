package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/tnynlabs/wyrm/pkg/storage/postgres"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file (error: %v)", err)
	}

	db, err := postgres.GetFromEnv()
	if err != nil {
		log.Fatalf("Error loading db instance (error: %v)", err)
	}

	schemaPath := os.Getenv("DB_SCHEMA_PATH")
	schemaGenQuery, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("Error loading schema file (%s) (error: %v)", schemaPath, err)
	}

	_, err = db.Exec(string(schemaGenQuery))
	if err != nil {
		log.Fatalf("Error initializing database schema (error: %v)", err)
	}

	log.Println("Database schema initialized successfully")
}
