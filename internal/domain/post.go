package domain

import "time"

type Post struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (*Post) NewPost(userID int64, title, content string) (*Post, error) {
	if title == "" || len(title) > 100 {
		return nil, ErrInvalidTitle

	}
	if content == "" || len(content) > 1000 || len(content) == 0 {
		return nil, ErrInvalidContent
	}
	return &Post{
		UserID:  userID,
		Title:   title,
		Content: content,
	}, nil
}
