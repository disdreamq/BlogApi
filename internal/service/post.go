package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"github.com/disdreamq/BlogApi/internal/domain"
	"github.com/disdreamq/BlogApi/internal/port"
	"github.com/rs/zerolog/log"
)

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
	logger := log.Ctx(ctx)
	trace_id, _ := ctx.Value("trace_id").(string)
	logger.Info().
		Str("trace_id", trace_id).
		Int64("user_id", userID).
		Str("title", title).
		Msg("Created post")
	return post, nil
}

func (p *PostService) GetByID(ctx context.Context, postID int64) (*domain.Post, error) {
	cachedPost, ok := p.cache.Get(ctx, "post_"+strconv.FormatInt(postID, 10))
	if !ok {
		post, err := p.postRepo.GetByID(ctx, postID)
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

		p.cache.Set(ctx, "post_"+strconv.FormatInt(postID, 10), data, 10*time.Minute)
		return post, nil
	}

	var post domain.Post
	err := json.Unmarshal([]byte(cachedPost), &post)
	if err != nil {
		return nil, ErrCacheUnmarshal
	}
	logger := log.Ctx(ctx)
	trace_id, _ := ctx.Value("trace_id").(string)
	logger.Debug().
		Str("trace_id", trace_id).
		Int64("post_id", postID).
		Msg("Read post")
	return &post, nil
}

func (p *PostService) GetByTitle(ctx context.Context, title string) (*domain.Post, error) {
	cachedPost, ok := p.cache.Get(ctx, "post_"+title)
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

		p.cache.Set(ctx, "post_"+title, data, 10*time.Minute)
		return post, nil

	}

	var post domain.Post
	err := json.Unmarshal([]byte(cachedPost), &post)
	if err != nil {
		return nil, ErrCacheUnmarshal
	}
	logger := log.Ctx(ctx)
	trace_id, _ := ctx.Value("trace_id").(string)
	logger.Debug().
		Str("trace_id", trace_id).
		Str("title", title).
		Msg("Read post")
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
	logger := log.Ctx(ctx)
	trace_id, _ := ctx.Value("trace_id").(string)
	logger.Debug().
		Str("trace_id", trace_id).
		Int64("post_id", postID).
		Msg("Updated post")
	p.cache.Del(ctx, "post_"+strconv.FormatInt(postID, 10))
	p.cache.Del(ctx, "post_"+title)
	return nil

}

func (p *PostService) Delete(ctx context.Context, currUserID int64, postID int64) error {
	if ok := p.validateCurrUser(ctx, currUserID, postID); !ok {
		return ErrMethodNotAllowed
	}
	title, err := p.postRepo.Delete(ctx, postID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrPostNotFound
		default:
			return ErrUnexpected
		}
	}
	logger := log.Ctx(ctx)
	trace_id, _ := ctx.Value("trace_id").(string)
	logger.Debug().
		Str("trace_id", trace_id).
		Int64("post_id", postID).
		Msg("Deleted post")
	p.cache.Del(ctx, "post_"+strconv.FormatInt(postID, 10))
	p.cache.Del(ctx, "post_"+title)
	return nil
}
func (p *PostService) validateCurrUser(ctx context.Context, currUserID int64, postID int64) bool {
	post, err := p.GetByID(ctx, postID)
	if err != nil {
		return false
	}
	logger := log.Ctx(ctx)
	trace_id, _ := ctx.Value("trace_id").(string)
	if post.UserID != currUserID {
		logger.Debug().
			Str("trace_id", trace_id).
			Int64("curr_user_id", currUserID).
			Int64("user_id", post.UserID).
			Msg("Validation failed for user.")
		return false
	}
	logger.Debug().
		Str("trace_id", trace_id).
		Int64("curr_user_id", currUserID).
		Int64("user_id", post.UserID).
		Msg("Validated user.")
	return true
}
