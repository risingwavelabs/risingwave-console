-- name: CreateDatabaseConnection :one
INSERT INTO database_connections (
    name,
    cluster_id,
    username,
    password,
    database,
    organization_id
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetOrgDatabaseConnection :one
SELECT * FROM database_connections
WHERE id = $1 AND organization_id = $2;

-- name: GetDatabaseConnectionByID :one
SELECT * FROM database_connections
WHERE id = $1;

-- name: ListOrgDatabaseConnections :many
SELECT * FROM database_connections
WHERE organization_id = $1
ORDER BY name;

-- name: UpdateOrgDatabaseConnection :one
UPDATE database_connections
SET
    name = $3,
    cluster_id = $4,
    username = $5,
    password = $6,
    database = $7,
    organization_id = $8,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND organization_id = $2
RETURNING *;

-- name: DeleteOrgDatabaseConnection :exec
DELETE FROM database_connections
WHERE id = $1 AND organization_id = $2;

-- name: GetOrgDatabaseByID :one
SELECT * FROM database_connections
WHERE id = $1 AND organization_id = $2;

-- name: GetAllOrgDatabseConnectionsByClusterID :many
SELECT * FROM database_connections
WHERE cluster_id = $1 AND organization_id = $2;

-- name: DeleteAllOrgDatabaseConnectionsByClusterID :exec
DELETE FROM database_connections
WHERE cluster_id = $1 AND organization_id = $2;
