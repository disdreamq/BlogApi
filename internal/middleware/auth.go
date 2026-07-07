package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/disdreamq/BlogApi/internal/service"
)

type AuthMiddleware struct {
	tokenProvider service.TokenProvider
}

func NewAuthMiddleware(tokenProvider *service.TokenProvider) *AuthMiddleware {
	return &AuthMiddleware{
		tokenProvider: *tokenProvider,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
		}
		token := parts[1]
		payload, err := m.tokenProvider.ValidateToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
		}
		ctx := context.WithValue(r.Context(), "user_id", payload.UserID)
		ctx = context.WithValue(ctx, "email", payload.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
