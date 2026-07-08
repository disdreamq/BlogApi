package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"github.com/disdreamq/BlogApi/internal/domain"
	"github.com/disdreamq/BlogApi/internal/port"
	"github.com/redis/go-redis/v9"
)

// TODO Добавить логгирование
type PostService struct {
	postRepo port.PostRepository
	cache    port.Cache
}

func (p *PostService) CreatePost(ctx context.Context, userID int64, title, content string) (*domain.Post, error) {
	domainPost, err := domain.NewPost(userID, title, content)
	if err != nil {
		return nil, err
	}
	post, err := p.postRepo.CreatePost(ctx, domainPost)
	if err != nil {
		return nil, ErrLinkedUserNotFound
	}
	return post, nil
}

func (p *PostService) GetPost(ctx context.Context, postID int64) (*domain.Post, error) {
	cachedPost, err := p.cache.Get(ctx, strconv.FormatInt(postID, 10))
	if err != nil {
		switch err {
		case redis.Nil:
			post, err := p.postRepo.ReadPost(ctx, postID)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
					return nil, ErrPostNotFound
				default:
					return nil, ErrUnexpected
				}
			}
			data, err := json.Marshal(post)
			if err != nil {
				return nil, err
			}

			p.cache.Set(ctx, string(rune(postID)), data, 10*time.Minute)
			return post, nil
		default:
			return nil, ErrUnexpected
		}
	}
	var post domain.Post
	err = json.Unmarshal(cachedPost, &post)
	if err != nil {
		return nil, ErrCacheUnmarshal
	}
	return &post, nil
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
	p.cache.Del(ctx, string(rune(post.ID)))
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
	p.cache.Del(ctx, string(rune(postID)))
	return nil
}
