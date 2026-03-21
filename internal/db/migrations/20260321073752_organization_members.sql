-- +goose Up

CREATE TYPE role AS ENUM ('admin', 'member', 'owner');

CREATE TABLE IF NOT EXISTS organization_members (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    role role NOT NULL DEFAULT 'member',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, organization_id)
);

-- +goose Down
DROP TABLE IF EXISTS organization_members;
DROP TYPE IF EXISTS role;
