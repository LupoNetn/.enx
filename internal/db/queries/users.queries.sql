-- name: CreateUser :one
INSERT INTO users (email, name, password) VALUES ($1, $2, $3) RETURNING id, email, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT id, email, name, created_at, updated_at FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, name, created_at, updated_at FROM users
WHERE id = $1;

-- name: GetUserForAuth :one
SELECT id, email, name, password, created_at, updated_at FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users 
SET 
  email = COALESCE(sqlc.narg('email'), email),
  name = COALESCE(sqlc.narg('name'),name),
  password = COALESCE(sqlc.narg('password'), password),
  updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING id, email, name, created_at, updated_at;


-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;


-- name: GetAllUserOrganizations :many
SELECT o.id,o.name,o.created_by,o.passkey,om.role FROM organizations o
INNER JOIN organization_members om ON om.organization_id = o.id
WHERE om.user_id = $1;
