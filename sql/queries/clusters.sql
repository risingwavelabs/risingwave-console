-- name: CreateCluster :one
INSERT INTO clusters (
    organization_id,
    name,
    host,
    sql_port,
    meta_port
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: InitCluster :one
INSERT INTO clusters (
    organization_id,
    name,
    host,
    sql_port,
    meta_port
) VALUES (
    $1, $2, $3, $4, $5
) ON CONFLICT (organization_id, name) DO UPDATE SET updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: GetOrgCluster :one
SELECT * FROM clusters
WHERE id = $1 AND organization_id = $2;

-- name: ListOrgClusters :many
SELECT * FROM clusters
WHERE organization_id = $1
ORDER BY name;

-- name: UpdateOrgCluster :one
UPDATE clusters
SET
    name = $3,
    host = $4,
    sql_port = $5,
    meta_port = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND organization_id = $2
RETURNING *;

-- name: DeleteOrgCluster :exec
DELETE FROM clusters
WHERE id = $1 AND organization_id = $2;

-- name: GetClusterByID :one
SELECT * FROM clusters
WHERE id = $1;
