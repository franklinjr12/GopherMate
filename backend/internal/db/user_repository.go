package db

import (
	"database/sql"
	"errors"
	"fmt"
	"gophermatebackend/internal/model"
	"gophermatebackend/internal/utils"
)

func CreateUser(user *model.User) error {
	db, err := InitDB()
	if err != nil {
		return err
	}

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

func GetUserByUsername(username string) (*model.User, error) {
	db, err := InitDB()
	if err != nil {
		return nil, err
	}

	var user model.User
	query := "SELECT id, username, password_hash FROM users WHERE username = $1"
	row := db.QueryRow(query, username)
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}
