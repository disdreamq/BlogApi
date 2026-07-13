package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"github.com/disdreamq/BlogApi/internal/domain"
	"github.com/disdreamq/BlogApi/internal/port"
)

// TODO F7. Получение списка постов автора по username с пагинацией (limit/offset), сортировка по created_at DESC (публично).
type PostService struct {
	postRepo port.PostRepository
	cache    port.Cache
}

func NewPostService(postRepo port.PostRepository, cache port.Cache) *PostService {
	return &PostService{postRepo: postRepo, cache: cache}
}

func (p *PostService) Create(ctx context.Context, userID int64, title, content string) (*domain.Post, error) {
	domainPost, err := domain.NewPost(userID, title, content)
	if err != nil {
		return nil, err
	}
	post, err := p.postRepo.Create(ctx, domainPost)
	if err != nil {
		return nil, ErrLinkedUserNotFound
	}
	return post, nil
}

func (p *PostService) GetByID(ctx context.Context, postID int64) (*domain.Post, error) {
	cachedPost, ok := p.cache.Get(ctx, "postTitle_"+strconv.FormatInt(postID, 10))
	if !ok {
		post, err := p.postRepo.GetByTitle(ctx, strconv.FormatInt(postID, 10))
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

		if ok := p.cache.Set(ctx, "postTitle_"+strconv.FormatInt(postID, 10), data, 10*time.Minute); !ok {
			return post, nil
		}
	}

	var post domain.Post
	err := json.Unmarshal([]byte(cachedPost), &post)
	if err != nil {
		return nil, ErrCacheUnmarshal
	}
	return &post, nil
}

func (p *PostService) GetByTitle(ctx context.Context, title string) (*domain.Post, error) {
	cachedPost, ok := p.cache.Get(ctx, "postTitle_"+title)
	if !ok {
		post, err := p.postRepo.GetByTitle(ctx, title)
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

		if ok := p.cache.Set(ctx, "postTitle_"+title, data, 10*time.Minute); !ok {
			return post, nil
		}
	}

	var post domain.Post
	err := json.Unmarshal([]byte(cachedPost), &post)
	if err != nil {
		return nil, ErrCacheUnmarshal
	}
	return &post, nil

}

func (p *PostService) Update(ctx context.Context, currUserID, postID int64, title, content string) error {
	if ok := p.validateCurrUser(ctx, currUserID, postID); !ok {
		return ErrMethodNotAllowed
	}
	post, err := domain.NewPost(currUserID, title, content)
	if err != nil {
		return err
	}
	err = p.postRepo.Update(ctx, post)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrPostNotFound
		default:
			return ErrUnexpected
		}
	}
	p.cache.Del(ctx, strconv.FormatInt(postID, 10))
	return nil

}

func (p *PostService) Delete(ctx context.Context, currUserID int64, postID int64) error {
	if ok := p.validateCurrUser(ctx, currUserID, postID); !ok {
		return ErrMethodNotAllowed
	}
	err := p.postRepo.Delete(ctx, postID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrPostNotFound
		default:
			return ErrUnexpected
		}
	}
	p.cache.Del(ctx, strconv.FormatInt(postID, 10))
	return nil
}
func (p *PostService) validateCurrUser(ctx context.Context, currUserID int64, postID int64) bool {
	post, err := p.GetByID(ctx, postID)
	if err != nil {
		return false
	}
	if post.UserID != currUserID {
		return false
	}
	return true
}
