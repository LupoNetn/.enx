-- name: CreateOrganization :one
INSERT INTO organizations (name,email,passkey,created_by) 
VALUES ($1,$2,$3,$4)
RETURNING id,name,email;

-- name: GetOrganizationByEmail :one
SELECT id,name,email,created_by FROM organizations
WHERE email = $1;

-- name: GetOrganizationByID :one
SELECT id,name,email,created_by FROM organizations
WHERE id = $1;

-- name: GetOrganizationOwner :one
SELECT u.id,u.email,u.name FROM users u
INNER JOIN organizations o ON o.created_by = u.id
WHERE o.id = $1;

-- name: UpdateOrganization :one
UPDATE organizations
SET 
  name = COALESCE(sqlc.narg('name'), name),
  email = COALESCE(sqlc.narg('email'), email),
  passkey = COALESCE(sqlc.narg('passkey'), passkey),
  updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING id,name,email;

-- name: DeleteOrganization :exec
DELETE FROM organizations WHERE id = $1;

-- name: AddUserToOrganization :one
INSERT INTO organization_members (user_id,organization_id,role)
VALUES ($1,$2,$3)
RETURNING *;

-- name: DeleteUserFromOrganization :exec
DELETE FROM organization_members 
WHERE user_id = $1 AND organization_id = $2;

-- name: UpdateUserInOrganization :one
UPDATE organization_members
SET 
  role = COALESCE(sqlc.narg('role'), role)
WHERE user_id = $1 AND organization_id = $2
RETURNING role;
        
-- name: GetAllUsersInOrganization :many
SELECT u.id,u.name,u.email FROM users u 
INNER JOIN organization_members om ON om.user_id = u.id
WHERE om.organization_id = $1;

-- name: GetOrganizationByName :one
SELECT id,name,email,created_by FROM organizations
WHERE name = $1;

-- name: GetAllOrganizationsByUser :many
SELECT o.id,o.name,o.email,o.passkey FROM organizations o 
INNER JOIN organization_members om ON om.organization_id = o.id
WHERE om.user_id = $1;
