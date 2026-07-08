package postgres

import "errors"

var (
	ErrNoRows = errors.New("No rows affected")
)
