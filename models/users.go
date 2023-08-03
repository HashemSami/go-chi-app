package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

type NewUser struct {
	Email    string
	Password string
}

func (us *UserService) Create(nu NewUser) (*User, error) {
	email := strings.ToLower(nu.Email)

	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(nu.Password), bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	passwordHash := string(hashedBytes)

	row := us.DB.QueryRow(
		`INSERT INTO users(email, password_hash)
	VALUES($1, $2) RETURNING id`,
		email, passwordHash)

	user := User{
		Email:        email,
		PasswordHash: string(hashedBytes),
	}

	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost,
	)
	if err != nil {
		fmt.Printf("error hashing: %v\n", password)
		return "", fmt.Errorf("create user: %w", err)
	}

	return string(hashedBytes), nil
}
