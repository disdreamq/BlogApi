-- +migrate NoTransaction
CREATE INDEX CONCURRENTLY idx_users_email ON users (email);
CREATE INDEX CONCURRENTLY idx_posts_title_user_id ON posts (title, user_id);