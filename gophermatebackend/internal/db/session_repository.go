package db

import (
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
