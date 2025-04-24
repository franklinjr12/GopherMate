package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	connStr := "user=postgres password=yourpassword dbname=gophermate sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, err
	}

	log.Println("Database connection established")
	return db, nil
}
