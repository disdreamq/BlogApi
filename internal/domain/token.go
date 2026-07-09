package domain

import "time"

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
