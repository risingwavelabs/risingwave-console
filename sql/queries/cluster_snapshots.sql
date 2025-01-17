-- name: CreateClusterSnapshot :exec
INSERT INTO cluster_snapshots (cluster_id, snapshot_id, name)
VALUES ($1, $2, $3);

-- name: DeleteClusterSnapshot :exec
DELETE FROM cluster_snapshots
WHERE cluster_id = $1 AND snapshot_id = $2;

-- name: ListClusterSnapshots :many
SELECT * FROM cluster_snapshots
WHERE cluster_id = $1 ORDER BY created_at DESC;
