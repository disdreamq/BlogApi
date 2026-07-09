package service

import (
	"context"
	"time"

	"github.com/disdreamq/BlogApi/internal/port"
)

type TokenProvider interface {
	GenerateToken(ctx context.Context, userID int64, email string) (string, error)
	ValidateToken(tokenString string) (*TokenPayload, error)
	RefreshToken(oldToken string) (string, error)
}

type TokenPayload struct {
	UserID   int64     `json:"user_id"`
	Email    string    `json:"email"`
	ExpireAt time.Time `json:"expire_at"`
}

type AuthResult struct {
	Token        string `json:"token"`
	TokenPayload *TokenPayload
}

func NewAuthResult(token string, tp *TokenPayload) *AuthResult {
	return &AuthResult{
		Token:        token,
		TokenPayload: tp,
	}
}

type AuthService struct {
	userService   port.UserRepository
	hasher        port.Hasher
	tokenProvider TokenProvider
	tokenTTL      time.Duration
}

func NewAuthService(
	userService port.UserRepository,
	hasher port.Hasher,
	tokenProvider TokenProvider,
	tokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		userService:   userService,
		hasher:        hasher,
		tokenProvider: tokenProvider,
		tokenTTL:      tokenTTL,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
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
	return NewAuthResult(token, payload), nil
}
