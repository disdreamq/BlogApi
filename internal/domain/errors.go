package domain

import "errors"

var (
	// User
	ErrInvalidUserName = errors.New("Invalid user name")
	ErrInvalidEmail    = errors.New("Invalid email")
	ErrInvalidPassword = errors.New("Invalid password")

	// Post
	ErrInvalidTitle   = errors.New("Invalid title")
	ErrInvalidContent = errors.New("Invalid content")
)
