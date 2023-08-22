package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailTaken = errors.New("models: email address is already in use")

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type NewUser struct {
	Email    string
	Password string
}

type UserService struct {
	DB *sql.DB
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

	user := User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	row := us.DB.QueryRow(
		`INSERT INTO users(email, password_hash)
	VALUES($1, $2) RETURNING id`,
		email, passwordHash)

	err = row.Scan(&user.ID)
	if err != nil {
		var pgError *pgconn.PgError
		// here we are passing a pointer to the PgError pointer
		// so can errors.As matching the two pointers
		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return nil, ErrEmailTaken
			}
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}

func (us *UserService) Authenticate(nu NewUser) (*User, error) {
	email := strings.ToLower(nu.Email)
	user := User{
		Email: email,
	}

	row := us.DB.QueryRow(
		`SELECT id, password_hash FROM users WHERE email=$1`,
		email,
	)

	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(nu.Password))
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	fmt.Println("Password is correct!!")

	return &user, nil
}

func (us *UserService) UpdatePassword(userID int, password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	passwordHash := string(hashedBytes)

	_, err = us.DB.Exec(`
		UPDATE users
		SET password_hash = $2
		WHERE id = $1;
	`, userID, passwordHash)

	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	return nil
}
