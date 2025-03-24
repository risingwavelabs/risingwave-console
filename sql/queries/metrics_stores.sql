-- name: GetMetricsStore :one
SELECT ms.*
FROM metrics_stores ms
    JOIN clusters c ON c.metrics_store_id = ms.id
WHERE c.id = $1;
