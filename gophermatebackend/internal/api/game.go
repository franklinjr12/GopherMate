package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"gophermatebackend/internal/db"
	"gophermatebackend/internal/utils"
)

func GamesHandler(w http.ResponseWriter, r *http.Request) {
	dbConn, err := db.InitDB()
	if err != nil {
		log.Printf("GamesHandler: Failed to initialize database: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	defer dbConn.Close()

	games, err := db.GetOpenGames(dbConn)
	if err != nil {
		log.Printf("GamesHandler: Failed to fetch games: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch games"})
		return
	}
	response := make([]map[string]interface{}, len(games))
	for i, game := range games {
		response[i] = map[string]interface{}{
			"id":           game.ID,
			"player_white": game.PlayerWhite.Int64,
			"player_black": game.PlayerBlack.Int64,
		}
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func JoinGameHandler(w http.ResponseWriter, r *http.Request) {
	dbConn, err := db.InitDB()
	if err != nil {
		log.Printf("JoinGameHandler: Failed to initialize database: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	defer dbConn.Close()

	// Parse game ID from URL: /api/games/{id}/join
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid game join URL"})
		return
	}
	gameID := parts[3]

	// Get user ID from session (for now, use a placeholder or 2)
	userID := int64(2) // TODO: Replace with session extraction

	// Attempt to join the game as black
	err = db.JoinGameAsBlack(dbConn, gameID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Game not found"})
			return
		}
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Joined game successfully"})
}

func MoveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Move submitted successfully"})
}

// CreateGameHandler handles POST /api/games
func CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	dbConn, err := db.InitDB()
	if err != nil {
		log.Printf("CreateGameHandler: Failed to initialize database: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	defer dbConn.Close()

	// TODO: Get user ID from session (for now, use a placeholder or 1)
	playerWhiteID := int64(1)

	gameID, err := db.CreateGame(dbConn, playerWhiteID)
	if err != nil {
		log.Printf("CreateGameHandler: Failed to create game: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create game"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"id": gameID})
}
