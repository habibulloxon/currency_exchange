package models

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func OpenExistingDB() error {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found; using environment variables from the system")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		return errors.New("DB_PATH not found in environment variables, please specify it in the .env file")
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("Database connection established using", dbPath)
	return nil
}
