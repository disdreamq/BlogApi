package service

import (
	"context"

	"github.com/disdreamq/BlogApi/internal/domain"
	"github.com/disdreamq/BlogApi/internal/port"
)

type PostService struct {
	postRepo port.PostRepository
}

func (p *PostService) CreatePost(ctx context.Context, userID int64, title, content string) (int64, error) {
	domainPost, err := domain.NewPost(userID, title, content)
	if err != nil {
		return -1, err
	}
	id, err := p.postRepo.CreatePost(ctx, domainPost.UserID, domainPost.Title, domainPost.Content)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (p *PostService) GetPost(ctx context.Context, postID int64) (*domain.Post, error) {
	user, err := p.postRepo.ReadPost(ctx, postID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *PostService) UpdatePost(ctx context.Context, post *domain.Post) error {
	return p.postRepo.UpdatePost(ctx, post)

}

func (p *PostService) DeletePost(ctx context.Context, postID int64) error {
	return p.postRepo.DeletePost(ctx, postID)
}
