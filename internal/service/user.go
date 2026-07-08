package service

import (
	"context"
	"database/sql"

	"github.com/disdreamq/BlogApi/internal/domain"
	"github.com/disdreamq/BlogApi/internal/port"
)

// TODO логирование
type UserService struct {
	userRepo port.UserRepository
	hasher   port.Hasher
}

func (u *UserService) CreateUser(ctx context.Context, username, email, password string) (*domain.User, error) {
	if password == "" {
		return nil, ErrEmptyPassword
	}
	if len(password) < 8 || len(password) > 60 {
		return nil, ErrInvalidPasswordLength
	}
	passwordHash, err := u.hasher.Hash(password)
	if err != nil {
		return nil, ErrCanNotCalculatePassHash
	}
	domainUser, err := domain.NewUser(username, email, passwordHash)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepo.CreateUser(ctx, domainUser)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrUserAlreadyExists
		default:
			return nil, ErrUnexpected
		}

	}
	return user, nil
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
