-- name: CreateSnapshot :exec
INSERT INTO snapshots (cluster_id, snapshot_id)
VALUES ($1, $2);

-- name: ListSnapshots :many
SELECT * FROM snapshots
WHERE cluster_id = $1
ORDER BY created_at DESC;

-- name: DeleteSnapshot :exec
DELETE FROM snapshots
WHERE cluster_id = $1 AND snapshot_id = $2;
