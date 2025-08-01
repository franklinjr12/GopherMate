package db

import (
	"database/sql"

	"gophermatebackend/internal/utils"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() (*sql.DB, error) {
	if db != nil && db.Ping() == nil {
		return db, nil
	} else if db != nil {
		db.Close()
	}

	config := utils.LoadConfig()
	connStr := "user=" + config.DBUser + " password=" + config.DBPassword + " dbname=" + config.DBName + " host=" + config.DBHost + " port=" + config.DBPort + " sslmode=require"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		utils.LogError("Error opening database: " + err.Error())
		return nil, err
	}

	if err := db.Ping(); err != nil {
		utils.LogError("Error connecting to database: " + err.Error())
		return nil, err
	}

	return db, nil
}
