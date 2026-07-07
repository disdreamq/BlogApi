package service

import "errors"

var (
	// Registration
	ErrEmptyPassword           = errors.New("Empty password not allowed.")
	ErrInvalidPasswordLength   = errors.New("Invalid length of user password.")
	ErrUserAlreadyExists       = errors.New("User with this email already exists.")
	ErrCanNotCalculatePassHash = errors.New("Can not calculate password hash.")

	// User
	ErrUserNotFound = errors.New("User not found.")
	ErrUserTimeOut  = errors.New("Time out error while processing user.")

	// Post
	ErrPostNotFound       = errors.New("Post not found.")
	ErrLinkedUserNotFound = errors.New("Linked user not found.")

	// Unexpected
	ErrUnexpected = errors.New("Unexpected error.")
)
