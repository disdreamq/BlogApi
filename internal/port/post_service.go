package port

import (
	"context"

	"github.com/disdreamq/BlogApi/internal/domain"
)

type PostService interface {
	CreatePost(ctx context.Context, userID int64, title, content string) (int64, error)
	GetPost(ctx context.Context, postID int64) (*domain.Post, error)
	UpdatePost(ctx context.Context, post *domain.Post) error
	DeletePost(ctx context.Context, postID int64) error
}
