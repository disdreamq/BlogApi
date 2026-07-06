package service

import (
	"database/sql"
	"errors"

	"github.com/disdreamq/BlogApi/internal/domain"
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
		if errors.Is(err, sql.ErrNoRows) {
			return -1, ErrUserAlreadyExists
		}
	}
	return id, nil
}

func (u *UserService) GetUserByID(userID int64) (*domain.User, error) {
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) GetUserByEmail(email string) (*domain.User, error) {
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) UpdateUser(user *domain.User) error {
	return u.userRepo.UpdateUser(user)

}

func (u *UserService) DeleteUser(userID int64) error {
	return u.userRepo.DeleteUser(userID)
}
