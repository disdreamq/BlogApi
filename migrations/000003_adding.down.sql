-- +migrate NoTransaction
DROP INDEX CONCURRENTLY idx_users_email
DROP INDEX CONCURRENTLY idx_posts_title_user_id