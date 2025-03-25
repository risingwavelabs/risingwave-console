-- name: GetMetricsStore :one
SELECT ms.*
FROM metrics_stores ms
    JOIN clusters c ON c.metrics_store_id = ms.id
WHERE c.id = $1;

-- name: ListMetricsStoresByOrgID :many
SELECT * FROM metrics_stores
WHERE organization_id = $1;

-- name: CreateMetricsStore :one
INSERT INTO metrics_stores (name, spec, organization_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateMetricsStore :one
UPDATE metrics_stores
SET name = $2, 
    spec = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND organization_id = $4
RETURNING *;

-- name: DeleteMetricsStore :exec
DELETE FROM metrics_stores
WHERE id = $1 AND organization_id = $2;

-- name: GetMetricsStoreByIDAndOrgID :one
SELECT * FROM metrics_stores
WHERE id = $1 AND organization_id = $2;
