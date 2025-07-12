package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"gophermatebackend/internal/cache"
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

	// Parse session token from request body
	var req struct {
		PlayerToken string `json:"player_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	// Get user ID from session token
	userID, err := db.GetUserIDBySessionToken(dbConn, req.PlayerToken)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid user token"})
		return
	}

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
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	dbConn, err := db.InitDB()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	defer dbConn.Close()

	// Validate that the incoming User(token) from the data matches a existing user session from the database
	userID, err := db.GetUserIDBySessionToken(dbConn, moveReq.User)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid user token"})
		return
	}

	// Validate that the userID is from the game session provided by incoming Session
	ok, err := db.ValidateUserInGameSession(dbConn, moveReq.Session, userID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "User is not part of the game session"})
		return
	}

	// Determine which color the user is in this game
	color, err := db.GetUserColorInGame(dbConn, moveReq.Session, userID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to determine player color"})
		return
	}
	if color == "" {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "User is not a player in this game"})
		return
	}

	board := cache.GetBoard(moveReq.Session)
	if board == nil {
		utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Game not found"})
		return
	}

	// Determine whose turn it is using board.LastMove
	var turnColor string
	switch board.LastMove {
	case "white":
		turnColor = "black"
	case "black":
		turnColor = "white"
	default:
		turnColor = "white" // fallback to white if unset
	}
	if color != turnColor {
		utils.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "It is not your turn"})
		return
	}

	// Validate move
	valid, err := movevalidation.ValidateMove(board, movevalidation.MoveData{
		Piece: moveReq.Piece,
		From:  movevalidation.Position{Row: moveReq.From.Row, Col: moveReq.From.Col},
		To:    movevalidation.Position{Row: moveReq.To.Row, Col: moveReq.To.Col},
	})
	if err != nil {
		log.Printf("MoveHandler: Invalid move: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if !valid {
		log.Printf("MoveHandler: Invalid move: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid move"})
		return
	}

	// Build notation: color-piece e2->e4
	from := string(rune('a'+moveReq.From.Col)) + string(rune('1'+(7-moveReq.From.Row)))
	to := string(rune('a'+moveReq.To.Col)) + string(rune('1'+(7-moveReq.To.Row)))
	notation := moveReq.Piece + " " + from + "->" + to

	err = db.SaveMove(dbConn, moveReq.Session, userID, notation)
	if err != nil {
		log.Printf("MoveHandler: Failed to save move: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save move"})
		return
	}

	// Update board state in cache
	board.Squares[moveReq.To.Row][moveReq.To.Col] = moveReq.Piece
	board.Squares[moveReq.From.Row][moveReq.From.Col] = "" // Clear the from square
	if color == "white" {
		board.LastMove = "white"
	} else {
		board.LastMove = "black"
	}

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

	// Check if a board already exists for this gameID (should not, but for safety)
	// If not, create and cache a new board
	if cache.GetBoard(gameID) == nil {
		cache.SetBoard(gameID, cache.NewInitialBoard())
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"id": gameID})
}
