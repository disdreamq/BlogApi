package domain

import "time"

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewUser(userName, email, passwordHash string) (*User, error) {
	if userName == "" || len(userName) > 30 {
		return nil, ErrInvalidUserName
	}
	if email == "" {
		return nil, ErrInvalidEmail
	}

	return &User{
		Username:     userName,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}, nil
}
