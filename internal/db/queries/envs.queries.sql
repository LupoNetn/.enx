-- name: CreateEnv :one
INSERT INTO envs (name, project_id, variables, description, created_by)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetEnvByID :one
SELECT * FROM envs
WHERE id = $1;

-- name: GetEnvsByProject :many
SELECT * FROM envs
WHERE project_id = $1;

-- name: GetEnvByNameInProject :one
SELECT * FROM envs
WHERE name = $1 AND project_id = $2;

-- name: UpdateEnv :one
UPDATE envs
SET 
  name = COALESCE(sqlc.narg('name'), name),
  variables = COALESCE(sqlc.narg('variables'), variables),
  description = COALESCE(sqlc.narg('description'), description),
  version = version + 1,
  updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteEnv :exec
DELETE FROM envs
WHERE id = $1;
