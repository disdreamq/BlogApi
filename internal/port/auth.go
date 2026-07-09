package port

import (
	"context"

	"github.com/disdreamq/BlogApi/internal/domain"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (*domain.AuthResult, error)
}
