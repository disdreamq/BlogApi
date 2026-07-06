package service

import (
	"github.com/disdreamq/BlogApi/internal/port"
)

type UserService struct {
	userRepo port.UserRepository
	hasher   port.Hasher
}

func (u *UserService) CreateUser(username, email, password string) (int64, error) {
	passwordHash, err := u.hasher.Hash(password)
	if err != nil {
		return -1, nil
	}
	id, err := u.userRepo.CreateUser(username, email, passwordHash)
	if err != nil {
		return -1, err
	}
	return id, nil
}
