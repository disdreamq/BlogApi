package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/disdreamq/BlogApi/internal/domain"
	"github.com/disdreamq/BlogApi/internal/port"
)

type UserService struct {
	userRepo port.UserRepository
	hasher   port.Hasher
}

func (u *UserService) CreateUser(ctx context.Context, username, email, password string) (int64, error) {
	passwordHash, err := u.hasher.Hash(password)
	if err != nil {
		return -1, nil
	}
	domainUser, err := domain.NewUser(username, email, passwordHash)
	if err != nil {
		return -1, err
	}
	id, err := u.userRepo.CreateUser(ctx, domainUser.Username, domainUser.Email, domainUser.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, ErrUserAlreadyExists
		}
	}
	return id, nil
}

func (u *UserService) GetUserByID(ctx context.Context, userID int64) (*domain.User, error) {
	user, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) UpdateUser(ctx context.Context, user *domain.User) error {
	return u.userRepo.UpdateUser(ctx, user)

}

func (u *UserService) DeleteUser(ctx context.Context, userID int64) error {
	return u.userRepo.DeleteUser(ctx, userID)
}
