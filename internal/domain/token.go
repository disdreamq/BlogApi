package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type TokenPayload struct {
	Claims   Claims
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
func NewPayload(claims *Claims, expireAt time.Time) *TokenPayload {
	return &TokenPayload{
		Claims:   *claims,
		ExpireAt: expireAt,
	}
}
