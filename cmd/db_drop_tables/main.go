package main

import (
	"fmt"
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

	dropPath := os.Getenv("DB_DROP_SCHEMA_PATH")
	fmt.Println(dropPath)
	schemaGenQuery, err := ioutil.ReadFile(dropPath)
	if err != nil {
		log.Fatalf("Error loading schema file (%s) (error: %v)", dropPath, err)
	}

	_, err = db.Exec(string(schemaGenQuery))
	if err != nil {
		log.Fatalf("Error Dropping database schema (error: %v)", err)
	}

	log.Println("Database schema droped successfully")
}
