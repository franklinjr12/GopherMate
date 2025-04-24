package api

import (
	"encoding/json"
	"net/http"
)

func GamesHandler(w http.ResponseWriter, r *http.Request) {
	// Placeholder for listing all open games
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "List of open games"})
}

func JoinGameHandler(w http.ResponseWriter, r *http.Request) {
	// Placeholder for joining a game
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Joined game successfully"})
}

func MoveHandler(w http.ResponseWriter, r *http.Request) {
	// Placeholder for submitting a move
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Move submitted successfully"})
}
