package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/HashemSami/go-chi-app/rand"
)

const (
	// minimum number of bytes to be used in each token
	MinBytesPerToken = 32
)

type Session struct {
	ID        int
	UserID    int
	TokenHash string
	// Token is only set when creating a new session.
	// when look up a session, this will be left empty,
	// as we only store the hash of a session in our database
	// and we cannot reverse it into a raw token
	Token string
}

type SessionService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes
	// to use when generating each session token. If this value
	// is not set or is less than the MinBytesPerToken will be used.
	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	// 1. Create the session token
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	// 2. Create the session
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}

	// 2.1 Query for user's session
	// 2.2 if not found, create a new session for the user
	// 2.3 If found, update the user's session
	row := ss.DB.QueryRow(`
	INSERT INTO sessions(user_id, token_hash)
	VALUES ($1, $2)
	ON CONFLICT (user_id) DO
	UPDATE
	SET token_hash = $2
	RETURNING id;
		`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)

	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &session, nil
}

// func (ss *SessionService) User(token string) (*User, error) {
// 	// 1. Hash the session token
// 	tokenHash := ss.hash(token)

// 	// 2. Query for the session with that hash
// 	row := ss.DB.QueryRow(`
// 		SELECT user_id FROM sessions
// 		WHERE token_hash = $1;
// 		`, tokenHash)

// 	var user User
// 	err := row.Scan(&user.ID)
// 	if err != nil {
// 		return nil, fmt.Errorf("user: %w", err)
// 	}

// 	// 3. Using the userID from the session, we
// 	// need to query for the user
// 	row = ss.DB.QueryRow(`
// 	SELECT email, password_hash
// 	FROM users
// 	WHERE id=$1;
// 	`, user.ID)

// 	err = row.Scan(&user.Email, &user.PasswordHash)
// 	if err != nil {
// 		return nil, fmt.Errorf("user: %w", err)
// 	}

// 	// 4. We need to return the user
// 	return &user, nil
// }

func (ss *SessionService) User(token string) (*User, error) {
	// 1. Hash the session token
	tokenHash := ss.hash(token)

	row := ss.DB.QueryRow(`
		SELECT sessions.user_id, users.email, users.password_hash
		FROM sessions
		LEFT JOIN users on users.id = sessions.user_id
		WHERE token_hash = $1;
		`, tokenHash)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	return &user, nil
}

func (ss SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)
	// using DB.Exec because we don't care about returning
	// data after deleting the row
	_, err := ss.DB.Exec(`
	DELETE FROM sessions
	WHERE token_hash = $1;
	`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
