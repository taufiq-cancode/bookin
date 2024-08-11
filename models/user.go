package models

import (
	"errors"
	"log"
	"booking.com/db"
	"booking.com/utils"
)

var ErrEmailExists = errors.New("Email already exists")

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	var existingUser User
	err := db.DB.QueryRow("SELECT id FROM users WHERE email = ?", u.Email).Scan(&existingUser.ID)
	if err == nil && existingUser.ID > 0 {
		log.Printf("Email already exists: %v", u.Email)
		return ErrEmailExists
	}

	query := "INSERT INTO users(email, password) VALUES (?, ?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		log.Printf("Error preparing query: %v", err)
		return err
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPassword)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return err
	}

	userId, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return err
	}

	u.ID = userId
	return nil
}

func (u *User) ValidateCredentials() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)

	if err != nil {
		return errors.New("Credentials invalid")
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)

	if !passwordIsValid {
		return errors.New("Credentials invalid")
	}

	return nil
}
