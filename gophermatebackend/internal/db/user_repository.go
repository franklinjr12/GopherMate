package db

import (
	"errors"
	"gophermatebackend/internal/model"
	"gophermatebackend/internal/utils"
)

func CreateUser(user *model.User) error {
	db, err := InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`
	_, err = db.Exec(query, user.Username, user.Email, hashedPassword)
	if err != nil {
		return errors.New("failed to insert user into database")
	}

	return nil
}
