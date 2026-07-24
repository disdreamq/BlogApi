package domain

import "errors"

var (
	// User
	ErrInvalidUserName = errors.New("Invalid user name")
	ErrInvalidEmail    = errors.New("Invalid email")

	// Post
	ErrInvalidTitle   = errors.New("Invalid title")
	ErrInvalidContent = errors.New("Invalid content")
	ErrInvalidUserId  = errors.New("Invalid user ID")
	ErrInvalidID      = errors.New("Invalid ID")
)
