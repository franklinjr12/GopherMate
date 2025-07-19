package api

import (
	"encoding/json"
	"gophermatebackend/internal/cache"
	"gophermatebackend/internal/db"
	"gophermatebackend/internal/utils"
	"net/http"
	"strings"
)

// BoardStateHandler handles GET /api/games/{id}/board
func BoardStateHandler(w http.ResponseWriter, r *http.Request) {
	// Parse game ID from URL: /api/games/{id}/board
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid board URL"})
		return
	}
	gameID := parts[3]

	// Authenticate user from Authorization header (Bearer <token>)
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Missing or invalid Authorization header"})
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	dbConn, err := db.InitDB()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	defer dbConn.Close()

	userID, err := db.GetUserIDBySessionToken(dbConn, token)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid user token"})
		return
	}

	// Check user is part of the game (optional, for security)
	ok, err := db.ValidateUserInGameSession(dbConn, gameID, userID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	if !ok {
		utils.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "User is not part of the game session"})
		return
	}

	// Get last move for this game using db.GetLastMove
	moveNumber, notation, err := db.GetLastMove(dbConn, gameID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get last move"})
		return
	}
	resp := map[string]interface{}{
		"number":   moveNumber,
		"notation": notation, // format is: white-pawn e2->e4
	}
	board := cache.GetBoard(gameID)
	if board != nil && board.DrawOfferPending {
		resp["draw_offer"] = board.DrawOffer
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
