package service

import (
	"context"

	"github.com/disdreamq/BlogApi/internal/domain"
	"github.com/disdreamq/BlogApi/internal/port"
	"github.com/rs/zerolog/log"
)

type AuthService struct {
	userService   port.UserService
	hasher        port.Hasher
	tokenProvider port.TokenProvider
}

func NewAuthService(
	userService port.UserService,
	hasher port.Hasher,
	tokenProvider port.TokenProvider,
) *AuthService {
	return &AuthService{
		userService:   userService,
		hasher:        hasher,
		tokenProvider: tokenProvider,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.AuthResult, error) {
	logger := log.Ctx(ctx)
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
	logger.Info().
		Str("trace_id", ctx.Value("trace_id").(string)).
		Int64("user_id", user.ID).
		Msg("User loggined")
	return domain.NewAuthResult(token, payload), nil
}
