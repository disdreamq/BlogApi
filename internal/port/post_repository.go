package port

import (
	"github.com/disdreamq/BlogApi/internal/domain"
)

type PostReader interface {
	ReadPost(ID int64) (*domain.Post, error)
}

type PostCreater interface {
	CreatePost(post *domain.Post) error
}

type PostUpdater interface {
	UpdatePost(post *domain.Post) error
}
type PostDeleter interface {
	DeletePost(id int64) error
}

type PostRepository interface {
	PostReader
	PostCreater
	PostUpdater
	PostDeleter
}
