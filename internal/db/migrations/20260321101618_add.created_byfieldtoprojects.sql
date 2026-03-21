-- +goose Up
ALTER TABLE projects
ADD COLUMN created_by UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL;

-- +goose Down
DROP COLUMN IF EXISTS created_by
