package port

import (
	"context"

	"github.com/disdreamq/BlogApi/internal/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, username, email, password string) (*domain.User, error)
	GetUserByID(ctx context.Context, userID int64) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, userID int64) error
}
