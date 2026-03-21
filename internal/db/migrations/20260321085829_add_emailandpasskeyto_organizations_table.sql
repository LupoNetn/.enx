-- +goose Up
ALTER TABLE organizations
ADD COLUMN email TEXT NOT NULL UNIQUE,
ADD COLUMN passkey TEXT NOT NULL;

-- +goose Down
DROP COLUMN IF EXISTS email;
DROP COLUMN IF EXISTS passkey;
