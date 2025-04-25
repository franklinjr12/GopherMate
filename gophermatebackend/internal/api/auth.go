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

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Placeholder for user login logic
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User logged in successfully"})
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
