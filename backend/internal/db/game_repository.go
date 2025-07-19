package db

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
)

// GetLastMove returns the last move number and notation for a game, or 0 and "" if none.
func GetLastMove(db *sql.DB, gameID string) (int, string, error) {
	row := db.QueryRow(`SELECT move_number, notation FROM moves WHERE game_id = $1 ORDER BY move_number DESC LIMIT 1`, gameID)
	var n sql.NullInt64
	var s sql.NullString
	err := row.Scan(&n, &s)
	if err != nil {
		// If no moves, return 0 and ""
		return 0, "", nil
	}
	moveNumber := 0
	notation := ""
	if n.Valid {
		moveNumber = int(n.Int64)
	}
	if s.Valid {
		notation = s.String
	}
	return moveNumber, notation, nil
}

// JoinGameAsBlack sets the player_black_id for a game if not already set.
func JoinGameAsBlack(db *sql.DB, gameID string, userID int64) error {
	// Only allow joining if player_black_id is NULL
	res, err := db.Exec(`UPDATE games SET player_black_id = $1 WHERE id = $2 AND player_black_id IS NULL`, userID, gameID)
	if err != nil {
		log.Printf("JoinGameAsBlack: Failed to update game: %v", err)
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

type Game struct {
	ID          string
	PlayerWhite sql.NullInt64
	PlayerBlack sql.NullInt64
	Winner      sql.NullString
	CreatedAt   string
	FinishedAt  sql.NullString
}

func GetOpenGames(db *sql.DB) ([]Game, error) {
	query := `SELECT id, player_white_id, player_black_id FROM games WHERE finished_at IS NULL`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("GetOpenGames: Failed to execute query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		var game Game
		if err := rows.Scan(&game.ID, &game.PlayerWhite, &game.PlayerBlack); err != nil {
			log.Printf("GetOpenGames: Failed to scan row: %v", err)
			return nil, err
		}
		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		log.Printf("GetOpenGames: Row iteration error: %v", err)
		return nil, err
	}

	return games, nil
}

// CreateGame inserts a new game into the database and returns the game ID.
func CreateGame(db *sql.DB, playerWhiteID int64) (string, error) {
	gameID := uuid.New().String()
	query := `INSERT INTO games (id, player_white_id) VALUES ($1, $2)`
	_, err := db.Exec(query, gameID, playerWhiteID)
	if err != nil {
		log.Printf("CreateGame: Failed to insert game: %v", err)
		return "", err
	}
	return gameID, nil
}

// ValidateUserInGameSession checks if the user is a participant in the game (white or black)
func ValidateUserInGameSession(db *sql.DB, gameID string, userID int64) (bool, error) {
	var count int
	query := `SELECT COUNT(1) FROM games WHERE id = $1 AND (player_white_id = $2 OR player_black_id = $2)`
	err := db.QueryRow(query, gameID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SetGameResigned sets the winner and finished_at for a game when a player resigns
func SetGameResigned(db *sql.DB, gameID string, winner string) error {
	query := `UPDATE games SET winner = $1, finished_at = NOW() WHERE id = $2`
	_, err := db.Exec(query, winner, gameID)
	if err != nil {
		log.Printf("SetGameResigned: Failed to update game: %v", err)
		return err
	}
	return nil
}
