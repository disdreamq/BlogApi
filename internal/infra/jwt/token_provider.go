package jwt

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type Provider struct {
	secret  []byte
	signing jwt.SigningMethod
	expiry  time.Duration
}

func NewProvider(secret string, expiry time.Duration) *Provider {
	return &Provider{
		secret:  []byte(secret),
		signing: jwt.SigningMethodHS256,
		expiry:  expiry,
	}
}

func (p *Provider) GenerateToken(_ context.Context, userID int64, email string) (string, error) {
	claims := Claims{
		UserId: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		}}
	token := jwt.NewWithClaims(p.signing, claims)
	return token.SignedString(p.secret)
}

func (p *Provider) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method != p.signing {
			return nil, ErrInvalidToken
		}
		return p.secret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(*Claims)
	if ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
