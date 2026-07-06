package port

import (
	"context"

	"github.com/disdreamq/BlogApi/internal/domain"
)

type UserReaderByID interface {
	GetUserByID(ctx context.Context, userID int64) (*domain.User, error)
}
type UserReaderByEmail interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type UserCreater interface {
	CreateUser(ctx context.Context, username, email, password string) (int64, error)
}

type UserUpdater interface {
	UpdateUser(ctx context.Context, user *domain.User) error
}
type UserDeleter interface {
	DeleteUser(ctx context.Context, id int64) error
}

type Hasher interface {
	Hash(password string) (string, error)
	Check(hashed, plain string) error
}

type UserRepository interface {
	UserReaderByID
	UserReaderByEmail
	UserCreater
	UserUpdater
	UserDeleter
}
