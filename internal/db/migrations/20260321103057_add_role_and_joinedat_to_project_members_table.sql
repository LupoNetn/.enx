-- +goose Up
ALTER TABLE project_members
ADD COLUMN role project_role NOT NULL DEFAULT 'member',
ADD COLUMN joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- +goose Down
ALTER TABLE project_mambers
DROP COLUMN IF EXISTS role,
DROP COLUMN IF EXISTS joined_at;
