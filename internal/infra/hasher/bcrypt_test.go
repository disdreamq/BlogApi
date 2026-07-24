package hasher

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestNewBcryptHasher(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		h := NewBcryptHasher(10)
		if h == nil {
			t.Fatal("expected non-nil hasher")
		}
		if h.cost != 10 {
			t.Fatalf("wanted cost to be 10, got %d", h.cost)
		}
	})

	t.Run("less than min cost", func(t *testing.T) {
		h := NewBcryptHasher(1)
		if h == nil {
			t.Fatal("expected non-nil hasher")
		}
		if h.cost != bcrypt.DefaultCost {
			t.Fatalf("wanted cost to be 10, got %d", h.cost)
		}
	})
}

func TestBcryptHasher_Hash(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		h := NewBcryptHasher(10)
		passToHash := "password"
		hash, err := h.Hash(passToHash)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(hash) == 0 {
			t.Fatal("wanted not empty hash, got empty string")
		}
		if hash == passToHash {
			t.Fatal("wanted hash to be different from password")
		}
	})
}

func TestBcryptHasher_Check(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		h := NewBcryptHasher(10)
		passToHash := "password"
		hash, err := h.Hash(passToHash)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := h.Check(hash, passToHash); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
