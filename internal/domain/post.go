package domain

import "time"

type Post struct {
	ID        int64
	UserID    int64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPost(userID int64, title, content string) (*Post, error) {
	if title == "" || len(title) > 100 {
		return nil, ErrInvalidTitle
	}
	if content == "" || len(content) > 1000 {
		return nil, ErrInvalidContent
	}
	if userID <= 0 {
		return nil, ErrInvalidUserId
	}

	return &Post{
		UserID:  userID,
		Title:   title,
		Content: content,
	}, nil
}
