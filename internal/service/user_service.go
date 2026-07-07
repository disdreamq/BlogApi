package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/disdreamq/BlogApi/internal/domain"
	"github.com/disdreamq/BlogApi/internal/port"
)

// TODO логирование
type UserService struct {
	userRepo port.UserRepository
	hasher   port.Hasher
}

func (u *UserService) CreateUser(ctx context.Context, username, email, password string) (int64, error) {
	if password == "" {
		return -1, ErrEmptyPassword
	}
	if len(password) < 8 || len(password) > 60 {
		return -1, ErrInvalidPasswordLength
	}
	passwordHash, err := u.hasher.Hash(password)
	if err != nil {
		return -1, ErrCanNotCalculatePassHash
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
		switch err {
		case sql.ErrNoRows:
			return nil, ErrUserNotFound
		default:
			return nil, ErrUnexpected
		}
	}
	return user, nil
}

func (u *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrUserNotFound
		default:
			return nil, ErrUnexpected
		}
	}
	return user, nil
}

func (u *UserService) UpdateUser(ctx context.Context, user *domain.User) error {
	err := u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrUserNotFound
		case sql.ErrTxDone:
			return ErrUserTimeOut
		default:
			return ErrUnexpected
		}
	}
	return nil
}

func (u *UserService) DeleteUser(ctx context.Context, userID int64) error {
	err := u.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrUserNotFound
		case sql.ErrTxDone:
			return ErrUserTimeOut
		default:
			return ErrUnexpected
		}
	}
	return nil
}
