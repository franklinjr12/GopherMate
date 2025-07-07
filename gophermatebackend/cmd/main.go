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
	mux.HandleFunc("/api/games/", func(w http.ResponseWriter, r *http.Request) {
		// Handle /api/games/{id}/join for joining a game
		if r.Method == http.MethodPost && len(r.URL.Path) > len("/api/games/") && r.URL.Path[len(r.URL.Path)-5:] == "/join" {
			api.JoinGameHandler(w, r)
			return
		}
		// Fallback to existing handler
		if r.Method == http.MethodPost {
			api.CreateGameHandler(w, r)
			return
		}
		api.GamesHandler(w, r)
	})
	mux.HandleFunc("/api/games/move", api.MoveHandler)

	// Wrap the mux with the CORS middleware
	handler := api.CORSMiddleware(mux)

	// Start HTTP server
	log.Printf("Server is running on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
