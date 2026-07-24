package domain

import (
	"strings"
	"time"
)

type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

func NewUser(username, email, passwordHash string) (*User, error) {
	if username == "" || len(username) > 30 {
		return nil, ErrInvalidUserName
	}
	if email == "" || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return nil, ErrInvalidEmail
	}

	return &User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}, nil
}
