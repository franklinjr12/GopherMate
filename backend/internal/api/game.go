package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"gophermatebackend/internal/db"
	"gophermatebackend/internal/movevalidation"
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
	log.Printf("MoveHandler: Received request for %s", r.URL.Path)
	// Parse request body
	var moveReq struct {
		Session string `json:"session"`
		User    string `json:"user"`
		Piece   string `json:"piece"`
		From    struct {
			Row int `json:"row"`
			Col int `json:"col"`
		} `json:"from"`
		To struct {
			Row int `json:"row"`
			Col int `json:"col"`
		} `json:"to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&moveReq); err != nil {
		log.Printf("MoveHandler: Failed to decode request body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	dbConn, err := db.InitDB()
	if err != nil {
		log.Printf("MoveHandler: Failed to initialize database: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	defer dbConn.Close()

	// Validate that the incoming User(token) from the data matches a existing user session from the database
	userID, err := db.GetUserIDBySessionToken(dbConn, moveReq.User)
	if err != nil {
		log.Printf("MoveHandler: Invalid user token: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid user token"})
		return
	}

	// Validate that the userID is from the game session provided by incoming Session
	ok, err := db.ValidateUserInGameSession(dbConn, moveReq.Session, userID)
	if err != nil {
		log.Printf("MoveHandler: Failed to validate user in game session: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	if !ok {
		log.Printf("MoveHandler: User %d is not part of game session %s", userID, moveReq.Session)
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "User is not part of the game session"})
		return
	}

	// Validate move
	valid, err := movevalidation.ValidateMove(movevalidation.MoveData{
		Piece: moveReq.Piece,
		From:  movevalidation.Position{Row: moveReq.From.Row, Col: moveReq.From.Col},
		To:    movevalidation.Position{Row: moveReq.To.Row, Col: moveReq.To.Col},
	})
	if err != nil {
		log.Printf("MoveHandler: Move validation error: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if !valid {
		log.Printf("MoveHandler: Invalid move for user %d in session %s", userID, moveReq.Session)
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid move"})
		return
	}

	// TODO: Save move to DB, update game state, etc.

	log.Printf("MoveHandler: Move submitted successfully by user %d in session %s", userID, moveReq.Session)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Move submitted successfully"})
}

// CreateGameHandler handles POST /api/games
func CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	dbConn, err := db.InitDB()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	defer dbConn.Close()

	var req struct {
		PlayerToken string `json:"player_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	// Get user ID from session token
	playerWhiteID, err := db.GetUserIDBySessionToken(dbConn, req.PlayerToken)
	if err != nil || playerWhiteID <= 0 {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid player token"})
		return
	}

	gameID, err := db.CreateGame(dbConn, playerWhiteID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create game"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"id": gameID})
}
