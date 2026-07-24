package jwt

import (
	"context"
	"testing"
	"time"
)

func TestProvider_GenerateAndValidate(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		p := NewProvider("secret", 1*time.Hour)

		ctx := context.Background()
		tokenStr, err := p.GenerateToken(ctx, 42, "test@mail.com")
		if err != nil {
			t.Fatalf("generate error: %v", err)
		}

		payload, err := p.ValidateToken(tokenStr)
		if err != nil {
			t.Fatalf("validate error: %v", err)
		}
		if payload.Claims.UserID != 42 {
			t.Errorf("got userID %d, want 42", payload.Claims.UserID)
		}
		if payload.Claims.Email != "test@mail.com" {
			t.Errorf("got email %q, want %q", payload.Claims.Email, "test@mail.com")
		}
	})
}
