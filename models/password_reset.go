package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/HashemSami/go-chi-app/rand"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID     int
	UserID int
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
	// verify we have a valid email address for a user, and get that user's ID
	email = strings.ToLower(email)

	var userID int
	row := prs.DB.QueryRow(`
		SELECT id FROM users
		WHERE email = $1;
	`, email)

	err := row.Scan(&userID)
	if err != nil {
		// TODO: consider returning a specific error when the user does not
		// exist, to allow the user to know that the email is not in the database
		// so they don't need to go and check their reset link in their inbox
		return nil, fmt.Errorf("create password reset: %w", err)
	}

	// Create the session token
	bytesPerToken := prs.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create password reset: %w", err)
	}

	duration := prs.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}

	pwReset := PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: prs.hash(token),
		ExpiresAt: time.Now().Add(duration),
	}

	// insert the PasswordReset into the DB
	row = prs.DB.QueryRow(`
	INSERT INTO password_resets (user_id, token_hash, expires_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_id) DO
	UPDATE
	SET token_hash = $2, expires_at = $3
	RETURNING id;
		`, pwReset.UserID, pwReset.TokenHash, pwReset.ExpiresAt)

	err = row.Scan(&pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("create password reset: %w", err)
	}

	return &pwReset, nil
}

func (prs *PasswordResetService) Consume(token string) (*User, error) {
	// query the user by using the password_resets table
	tokenHash := prs.hash(token)
	var user User
	var pwReset PasswordReset
	row := prs.DB.QueryRow(`
		SELECT password_resets.id,
			password_resets.expires_at,
			users.id,
			users.email,
			users.password_hash
		FROM password_resets
		JOIN users ON users.id = password_resets.user_id
		WHERE token_hash = $1;
	`, tokenHash)

	err := row.Scan(&pwReset.ID,
		&pwReset.ExpiresAt,
		&user.ID,
		&user.Email,
		&user.PasswordHash,
	)
	if err != nil {
		return nil, fmt.Errorf("consume password reset: %w", err)
	}

	// if the user returned at this point, we want to make sure that
	// the reset token has a valid expiry date
	if time.Now().After(pwReset.ExpiresAt) {
		// this code will execute if the date has expired
		return nil, fmt.Errorf("token expired: %v", token)
	}

	// at this point we need to delete the password reset token from
	// the database so it can not be used later
	err = prs.delete(pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("consume password reset: %w", err)
	}

	return &user, nil
}

func (prs *PasswordResetService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (prs *PasswordResetService) delete(id int) error {
	_, err := prs.DB.Exec(`
	  DELETE FROM password_resets
		WHERE id = $1;
	`, id)
	if err != nil {
		return fmt.Errorf("delete pw reset: %w", err)
	}

	return nil
}
