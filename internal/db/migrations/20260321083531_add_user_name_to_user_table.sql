-- +goose Up
ALTER TABLE users
ADD COLUMN name TEXT NOT NULL DEFAULT 'anonymous';

-- +goose Down
DROP COLUMN name FROM users;