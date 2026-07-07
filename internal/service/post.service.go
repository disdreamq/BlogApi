package service

import (
	"context"
	"database/sql"

	"github.com/disdreamq/BlogApi/internal/domain"
	"github.com/disdreamq/BlogApi/internal/port"
)

// TODO логирование и кеш
type PostService struct {
	postRepo port.PostRepository
	cache    port.Cache
}

func (p *PostService) CreatePost(ctx context.Context, userID int64, title, content string) (int64, error) {
	domainPost, err := domain.NewPost(userID, title, content)
	if err != nil {
		return -1, err
	}
	id, err := p.postRepo.CreatePost(ctx, domainPost.UserID, domainPost.Title, domainPost.Content)
	if err != nil {
		return -1, ErrLinkedUserNotFound
	}
	return id, nil
}

func (p *PostService) GetPost(ctx context.Context, postID int64) (*domain.Post, error) {
	user, err := p.postRepo.ReadPost(ctx, postID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrPostNotFound
		default:
			return nil, ErrUnexpected
		}
	}
	return user, nil
}

func (p *PostService) UpdatePost(ctx context.Context, post *domain.Post) error {
	err := p.postRepo.UpdatePost(ctx, post)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrPostNotFound
		default:
			return ErrUnexpected
		}
	}
	return nil

}

func (p *PostService) DeletePost(ctx context.Context, postID int64) error {
	err := p.postRepo.DeletePost(ctx, postID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrPostNotFound
		default:
			return ErrUnexpected
		}
	}
	return nil
}
