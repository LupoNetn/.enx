-- +goose Up
ALTER TABLE projects
ADD COLUMN passkey TEXT NOT NULL;

-- +goose Down
DROP COLUMN passkey IF EXISTS;
