package postgres

import (
	"context"
	"time"

	"github.com/disdreamq/BlogApi/internal/domain"
	"github.com/jmoiron/sqlx"
)

type postRow struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r postRow) toDomain() *domain.Post {
	return &domain.Post{
		ID:        r.ID,
		UserID:    r.UserID,
		Title:     r.Title,
		Content:   r.Content,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

type PostRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(ctx context.Context, userID int64, title, content string) (int64, error) {

	tx, err := r.db.Beginx()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	txCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `
        INSERT INTO posts (user_id, title, content)
        VALUES ($1, $2, $3)
        RETURNING id
    `
	var id int64
	err = tx.GetContext(txCtx, &id, query, userID, title, content)
	if err != nil {
		return -1, err
	}
	if err = tx.Commit(); err != nil {
		return -1, err
	}

	return id, nil
}

func (r *PostRepository) ReadPost(ctx context.Context, userID int64) (*domain.Post, error) {
	txCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `
        SELECT * FROM posts
        WHERE id = $1
    `
	var post postRow
	err := r.db.GetContext(txCtx, &post, query, userID)
	if err != nil {
		return nil, err
	}
	return post.toDomain(), nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, post *domain.Post) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	txCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `
		UPDATE posts SET title = $1, content = $2
		WHERE id = $3
	`
	var postRow postRow
	err = tx.GetContext(txCtx, &postRow, query, post.Title, post.Content, post.ID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) DeletePost(ctx context.Context, postID int64) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	txCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `
		DELETE FROM posts WHERE id = $1
	`
	_, err = tx.ExecContext(txCtx, query, postID)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
