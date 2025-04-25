package db

import (
	"database/sql"
	"log"
)

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
