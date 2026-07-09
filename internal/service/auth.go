package service

import (
	"context"
	"time"

	"github.com/disdreamq/BlogApi/internal/domain"
	"github.com/disdreamq/BlogApi/internal/port"
)

type AuthService struct {
	userService   port.UserRepository
	hasher        port.Hasher
	tokenProvider port.TokenProvider
	tokenTTL      time.Duration
}

func NewAuthService(
	userService port.UserRepository,
	hasher port.Hasher,
	tokenProvider port.TokenProvider,
	tokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		userService:   userService,
		hasher:        hasher,
		tokenProvider: tokenProvider,
		tokenTTL:      tokenTTL,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.AuthResult, error) {
	user, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if err := s.hasher.Check(user.PasswordHash, password); err != nil {
		return nil, ErrWrongPassword
	}
	token, err := s.tokenProvider.GenerateToken(ctx, user.ID, user.Email)
	if err != nil {
		return nil, ErrCanNotLogin
	}
	payload, _ := s.tokenProvider.ValidateToken(token)
	return domain.NewAuthResult(token, payload), nil
}
