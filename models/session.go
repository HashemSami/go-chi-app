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

	// TODO: Store the session in the database
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	return nil, nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}