-- name: GetMetricsStore :one
SELECT ms.*
FROM metrics_stores ms
    JOIN clusters c ON c.metrics_store_id = ms.id
WHERE c.id = $1;

-- name: ListMetricsStoresByOrgID :many
SELECT * FROM metrics_stores
WHERE org_id = $1;

-- name: CreateMetricsStore :one
INSERT INTO metrics_stores (name, spec, org_id, default_labels)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: InitMetricsStore :one
INSERT INTO metrics_stores (name, spec, org_id, default_labels)
VALUES ($1, $2, $3, $4)
ON CONFLICT (org_id, name) DO UPDATE 
    SET 
        spec = EXCLUDED.spec,
        default_labels = EXCLUDED.default_labels,
        updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: UpdateMetricsStore :one
UPDATE metrics_stores
SET name = $2, 
    spec = $3,
    default_labels = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND org_id = $5
RETURNING *;

-- name: DeleteMetricsStore :exec
DELETE FROM metrics_stores
WHERE id = $1 AND org_id = $2;

-- name: GetMetricsStoreByIDAndOrgID :one
SELECT * FROM metrics_stores
WHERE id = $1 AND org_id = $2;
