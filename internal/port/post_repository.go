package port

import (
	"context"

	"github.com/disdreamq/BlogApi/internal/domain"
)

type PostCreater interface {
	CreatePost(ctx context.Context, userID int64, title, content string) (int64, error)
}
type PostReader interface {
	ReadPost(ctx context.Context, ID int64) (*domain.Post, error)
}

type PostUpdater interface {
	UpdatePost(ctx context.Context, post *domain.Post) error
}
type PostDeleter interface {
	DeletePost(ctx context.Context, id int64) error
}

type PostRepository interface {
	PostReader
	PostCreater
	PostUpdater
	PostDeleter
}
