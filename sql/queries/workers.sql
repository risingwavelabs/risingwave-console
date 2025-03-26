-- name: SendWorkerHeartbeat :exec
INSERT INTO workers (worker_name, last_heartbeat)
VALUES ($1, $2)
ON CONFLICT (worker_name) DO UPDATE SET last_heartbeat = $2;

-- name: PullTask :one
SELECT * FROM tasks
WHERE 
    (remaning IS NULL OR remaning > 0)
    AND (
        (worker_name IS NULL AND status = 'pending')
        OR
        (worker_name = $1 AND (status = 'running' OR status = 'pending'))
    )
    AND (
        started_at IS NULL OR started_at < NOW()
    )
ORDER BY created_at ASC
FOR UPDATE SKIP LOCKED
LIMIT 1;

-- name: LockTask :one
UPDATE tasks
SET 
    status = 'running',
    worker_name = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateTaskMetadata :one
UPDATE tasks
SET 
    status = $2, 
    remaining = $3, 
    started_at = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateTaskSpec :one
UPDATE tasks
SET spec = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: CreateTask :one
INSERT INTO tasks (worker_name, spec, status, remaining, started_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: InsertEvent :one
INSERT INTO events (spec)
VALUES ($1)
RETURNING *;
