-- name: CreateProject :one
INSERT INTO projects (name, passkey, organization_id, created_by)
VALUES ($1, $2, $3, $4)
RETURNING id, name, organization_id, created_at;

-- name: GetProjectByID :one
SELECT id, name, organization_id, created_by, created_at
FROM projects
WHERE id = $1;

-- name: GetProjectByName :one
SELECT id, name, organization_id, created_by, created_at
FROM projects
WHERE name = $1 AND organization_id = $2;

-- name: GetProjectsByOrganization :many
SELECT id, name, organization_id, created_by, created_at
FROM projects
WHERE organization_id = $1;

-- name: GetProjectsByUser :many
SELECT p.id, p.name, p.organization_id, p.created_by, p.created_at, pm.role
FROM projects p
INNER JOIN project_members pm ON pm.project_id = p.id
WHERE pm.user_id = $1;

-- name: GetProjectOwner :one
SELECT u.id, u.email, u.name FROM users u
INNER JOIN projects p ON p.created_by = u.id
WHERE p.id = $1;

-- name: UpdateProject :one
UPDATE projects
SET
    name       = COALESCE(sqlc.narg('name'), name),
    passkey    = COALESCE(sqlc.narg('passkey'), passkey),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING id, name, organization_id, updated_at;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;

-- name: AddUserToProject :one
INSERT INTO project_members (user_id, project_id, role)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAllUsersInProject :many
SELECT u.id, u.name, u.email, pm.role, pm.joined_at
FROM users u
INNER JOIN project_members pm ON pm.user_id = u.id
WHERE pm.project_id = $1;

-- name: GetUserRoleInProject :one
SELECT role FROM project_members
WHERE user_id = $1 AND project_id = $2;

-- name: UpdateUserInProject :one
UPDATE project_members
SET role = COALESCE(sqlc.narg('role'), role)
WHERE user_id = $1 AND project_id = $2
RETURNING role;

-- name: DeleteUserFromProject :exec
DELETE FROM project_members
WHERE user_id = $1 AND project_id = $2;