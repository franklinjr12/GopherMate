package db

import (
	"database/sql"
	"fmt"
	"log"

	"gophermatebackend/internal/utils"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	config := utils.LoadConfig()
	connStr := "user=" + config.DBUser + " password=" + config.DBPassword + " dbname=" + config.DBName + " host=" + config.DBHost + " port=" + config.DBPort + " sslmode=disable"
	fmt.Println("Using connection string:", connStr)
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
