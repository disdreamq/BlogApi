package domain

import "time"

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

func (u *User) NewUser(userName, email, password string) (*User, error) {
	if userName == "" || len(userName) > 30 {
		return nil, ErrInvalidUserName
	}
	if email == "" {
		return nil, ErrInvalidEmail
	}
	if password == "" || len(password) < 8 || len(password) > 30 {
		return nil, ErrInvalidPassword
	}

	return &User{
		Username:     userName,
		Email:        email,
		PasswordHash: password,
		CreatedAt:    time.Now(),
	}, nil
}
