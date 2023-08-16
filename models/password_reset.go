package models

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID     int
	UserId int
	// Token is only set when a PasswordReset is being created
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes
	// to use when generating each session token. If this value
	// is not set or is less than the MinBytesPerToken will be used.
	BytesPerToken int
	// Duration is the amount of time that a PasswordReset is valid for
	// Defaults to DefaultResetDuration
	Duration time.Duration
}

func (prs *PasswordResetService) Create(email string) (*PasswordReset, error) {
	return nil, fmt.Errorf("TODO:implement PasswordResetService.Create")
}

func (prs *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: implement the Consume function")
}
