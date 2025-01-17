-- name: CreateClusterDiagnostic :one
INSERT INTO cluster_diagnostics (cluster_id, content) VALUES ($1, $2) RETURNING *;

-- name: GetClusterDiagnostic :one
SELECT * FROM cluster_diagnostics WHERE id = $1;

-- name: ListClusterDiagnostics :many
SELECT id, created_at FROM cluster_diagnostics WHERE cluster_id = $1 ORDER BY created_at DESC;
