package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func CreateSession(userID int) (string, error) {
	db, err := InitDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	sessionToken := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	query := "INSERT INTO sessions (token, user_id, expires_at) VALUES ($1, $2, $3)"
	_, err = db.Exec(query, sessionToken, userID, expiresAt)
	if err != nil {
		return "", err
	}

	return sessionToken, nil
}

func GetUserIDBySessionToken(db *sql.DB, sessionToken string) (int64, error) {
	var userID int64
	query := "SELECT user_id FROM sessions WHERE token = $1 AND expires_at > NOW()"
	row := db.QueryRow(query, sessionToken)
	if err := row.Scan(&userID); err != nil {
		return 0, err
	}
	return userID, nil
}
