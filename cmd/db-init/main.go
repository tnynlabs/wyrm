package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file (error: %v)", err)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	schemaPath := os.Getenv("DB_SCHEMA_PATH")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}

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
