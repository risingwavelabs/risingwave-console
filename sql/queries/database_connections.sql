-- name: CreateDatabaseConnection :one
INSERT INTO database_connections (
    name,
    cluster_id,
    username,
    password,
    organization_id
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetDatabaseConnection :one
SELECT * FROM database_connections
WHERE id = $1;

-- name: ListDatabaseConnections :many
SELECT * FROM database_connections
WHERE organization_id = $1
ORDER BY name;

-- name: UpdateDatabaseConnection :one
UPDATE database_connections
SET
    name = $2,
    cluster_id = $3,
    username = $4,
    password = $5,
    organization_id = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteDatabaseConnection :exec
DELETE FROM database_connections
WHERE id = $1;
