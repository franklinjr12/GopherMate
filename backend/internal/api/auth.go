package api

import (
	"encoding/json"
	"log"
	"net/http"

	"gophermatebackend/internal/db"
	"gophermatebackend/internal/model"
	"gophermatebackend/internal/utils"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		log.Println("RegisterHandler: Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("RegisterHandler: Failed to decode request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := db.CreateUser(&user); err != nil {
		log.Printf("RegisterHandler: Failed to create user: %v\n", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Fetch the user to get the ID (in case it's not set)
	createdUser, err := db.GetUserByUsername(user.Username)
	if err != nil {
		log.Printf("RegisterHandler: Failed to fetch created user: %v\n", err)
		http.Error(w, "Failed to fetch user after registration", http.StatusInternalServerError)
		return
	}

	sessionToken, err := db.CreateSession(createdUser.ID)
	if err != nil {
		log.Printf("RegisterHandler: Failed to create session: %v\n", err)
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "User registered successfully", "token": sessionToken})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("LoginHandler: Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Printf("LoginHandler: Failed to decode request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByUsername(credentials.Username)
	if err != nil {
		log.Printf("LoginHandler: User not found: %v\n", err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(credentials.Password, user.Password) {
		log.Println("LoginHandler: Invalid password")
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	sessionToken, err := db.CreateSession(user.ID)
	if err != nil {
		log.Printf("LoginHandler: Failed to create session: %v\n", err)
		http.Error(w, "Failed to login", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "User logged in successfully", "token": sessionToken})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Placeholder for user logout logic
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User logged out successfully"})
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	// Placeholder for fetching current user info
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User info fetched successfully"})
}
