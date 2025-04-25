package api

import (
	"encoding/json"
	"log"
	"net/http"

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
	// Placeholder for joining a game
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Joined game successfully"})
}

func MoveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Move submitted successfully"})
}
