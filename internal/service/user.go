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

func (u *UserService) Create(ctx context.Context, username, email, password string) (*domain.User, error) {
	passwordHash, err := processPassword(password, u.hasher)
	if err != nil {
		return nil, err
	}
	domainUser, err := domain.NewUser(username, email, passwordHash)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepo.Create(ctx, domainUser)
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

func (u *UserService) GetByID(ctx context.Context, userID int64) (*domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
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

func (u *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
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

func (u *UserService) Update(ctx context.Context, currUserID, userID int64, username, email, password string) error {
	if ok := u.validateCurrUser(ctx, currUserID, userID); !ok {
		return ErrMethodNotAllowed
	}
	passwordHash, err := processPassword(password, u.hasher)
	if err != nil {
		return err
	}
	domainUser, err := domain.NewUser(username, email, passwordHash)
	if err != nil {
		return err
	}
	err = u.userRepo.Update(ctx, domainUser)
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

func (u *UserService) Delete(ctx context.Context, currUserID int64, userID int64) error {
	if ok := u.validateCurrUser(ctx, currUserID, userID); !ok {
		return ErrMethodNotAllowed
	}
	err := u.userRepo.Delete(ctx, userID)
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
func (p *UserService) validateCurrUser(ctx context.Context, currUserID, userID int64) bool {
	user, err := p.GetByID(ctx, userID)
	if err != nil {
		return false
	}
	if user.ID != currUserID {
		return false
	}
	return true
}

func processPassword(pass string, hasher port.Hasher) (string, error) {
	if pass == "" {
		return "", ErrEmptyPassword
	}
	if len(pass) < 8 || len(pass) > 60 {
		return "", ErrInvalidPasswordLength
	}
	pass, err := hasher.Hash(pass)
	if err != nil {
		return "", ErrCanNotCalculatePassHash
	}
	return pass, nil
}
