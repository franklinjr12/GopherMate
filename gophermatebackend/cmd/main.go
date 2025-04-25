package main

import (
	"log"
	"net/http"
	"os"

	"gophermatebackend/internal/api"
)

func main() {
	// Load environment variables or default values
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set up routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/register", api.RegisterHandler)
	mux.HandleFunc("/api/login", api.LoginHandler)
	mux.HandleFunc("/api/logout", api.LogoutHandler)
	mux.HandleFunc("/api/me", api.MeHandler)
	mux.HandleFunc("/api/games", api.GamesHandler)
	mux.HandleFunc("/api/games/join", api.JoinGameHandler)
	mux.HandleFunc("/api/games/move", api.MoveHandler)

	// Wrap the mux with the CORS middleware
	handler := api.CORSMiddleware(mux)

	// Start HTTP server
	log.Printf("Server is running on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
