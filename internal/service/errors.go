package service

import "errors"

var (
	// Registration
	ErrUserAlreadyExists = errors.New("User with this email already exists")
)
