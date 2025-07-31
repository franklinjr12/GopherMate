package db

import (
	"database/sql"

	"gophermatebackend/internal/utils"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	config := utils.LoadConfig()
	connStr := "user=" + config.DBUser + " password=" + config.DBPassword + " dbname=" + config.DBName + " host=" + config.DBHost + " port=" + config.DBPort + " sslmode=require"
	utils.LogInfo("Connecting to database: " + config.DBHost + ":" + config.DBPort)
	db, err := sql.Open("postgres", connStr)
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
