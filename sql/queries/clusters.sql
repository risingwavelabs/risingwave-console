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

-- name: GetCluster :one
SELECT * FROM clusters
WHERE id = $1;

-- name: ListClusters :many
SELECT * FROM clusters
WHERE organization_id = $1
ORDER BY name;

-- name: UpdateCluster :one
UPDATE clusters
SET
    name = $2,
    host = $3,
    sql_port = $4,
    meta_port = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteCluster :exec
DELETE FROM clusters
WHERE id = $1;
