package db

import (
	"database/sql"
	"time"

	"gophermatebackend/internal/cache"

	"github.com/google/uuid"
)

func defaultExpirationTime() time.Time {
	return time.Now().Add(24 * time.Hour)
}

func CreateSession(userID int) (string, error) {
	db, err := InitDB()
	if err != nil {
		return "", err
	}

	sessionToken := uuid.New().String()
	expiresAt := defaultExpirationTime()

	query := "INSERT INTO sessions (token, user_id, expires_at) VALUES ($1, $2, $3)"
	_, err = db.Exec(query, sessionToken, userID, expiresAt)
	if err != nil {
		return "", err
	}

	return sessionToken, nil
}

// Refactored: now uses cache for session token lookup
func GetUserIDBySessionToken(db *sql.DB, sessionToken string) (int64, error) {
	// Check cache first
	if userID, ok := cache.GetUserIDByToken(sessionToken); ok {
		return userID, nil
	}
	// Not in cache or expired, query DB
	var userID int64
	var expiresAt time.Time
	query := "SELECT user_id, expires_at FROM sessions WHERE token = $1"
	row := db.QueryRow(query, sessionToken)
	if err := row.Scan(&userID, &expiresAt); err != nil {
		// On error, do not cache
		return 0, err
	}
	// if the session is expired, delete the session from the cache and create a new one
	if expiresAt.Before(time.Now()) {
		cache.DeleteUserIDForToken(sessionToken)
		var err error
		sessionToken, err = CreateSession(int(userID))
		if err != nil {
			return 0, err
		}
	}
	// Update cache
	cache.SetUserIDForToken(sessionToken, userID, defaultExpirationTime())
	return userID, nil
}
