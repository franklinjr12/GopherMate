package db

import (
	"database/sql"
	"fmt"
)

// SaveMove inserts a move into the moves table.
// SaveMove inserts a move into the moves table. move_number is set by DB trigger.
func SaveMove(dbConn *sql.DB, gameID string, playerID int64, notation string) error {
	query := `INSERT INTO moves (game_id, player_id, notation) VALUES ($1, $2, $3)`
	_, err := dbConn.Exec(query, gameID, playerID, notation)
	if err != nil {
		return fmt.Errorf("failed to save move: %w", err)
	}
	return nil
}

// No need for GetMoveCount; move_number is handled by DB trigger.
