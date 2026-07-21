package port

import (
	"context"

	"github.com/disdreamq/BlogApi/internal/domain"
)

type PostCreater interface {
	Create(ctx context.Context, post *domain.Post) (*domain.Post, error)
}
type PostReaderByID interface {
	GetByID(ctx context.Context, ID int64) (*domain.Post, error)
}
type PostReaderByTitle interface {
	GetByTitle(ctx context.Context, title string) (*domain.Post, error)
}

type PostUpdater interface {
	Update(ctx context.Context, post *domain.Post) error
}
type PostDeleter interface {
	Delete(ctx context.Context, id int64) (string, error)
}

type PostRepository interface {
	PostReaderByID
	PostReaderByTitle
	PostCreater
	PostUpdater
	PostDeleter
}
