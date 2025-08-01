package db

import (
	"database/sql"
)

// GetUserColorInGame returns "white" or "black" if the user is a player in the game, or "" if not
func GetUserColorInGame(db *sql.DB, gameID string, userID int64) (string, error) {
	var whiteID, blackID int64
	query := `SELECT player_white_id, player_black_id FROM games WHERE id = $1`
	err := db.QueryRow(query, gameID).Scan(&whiteID, &blackID)
	if err != nil {
		return "", err
	}
	if whiteID == userID {
		return "white", nil
	}
	if blackID == userID {
		return "black", nil
	}
	return "", nil
}
