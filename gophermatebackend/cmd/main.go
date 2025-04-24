package main

import (
	"log"
	"net/http"
	"os"

	"gophermatebackend/internal/api"
	"gophermatebackend/internal/db"
)

func main() {
	// Load environment variables or default values
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database connection
	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Set up routes
	http.HandleFunc("/api/register", api.RegisterHandler)
	http.HandleFunc("/api/login", api.LoginHandler)
	http.HandleFunc("/api/logout", api.LogoutHandler)
	http.HandleFunc("/api/me", api.MeHandler)
	http.HandleFunc("/api/games", api.GamesHandler)
	http.HandleFunc("/api/games/join", api.JoinGameHandler)
	http.HandleFunc("/api/games/move", api.MoveHandler)

	// Start HTTP server
	log.Printf("Server is running on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
