-- name: PullTask :one
SELECT * FROM tasks
WHERE 
    status = 'pending'
    AND (
        started_at IS NULL OR started_at < NOW()
    )
ORDER BY created_at ASC
FOR UPDATE SKIP LOCKED
LIMIT 1;

-- name: UpdateTaskStatus :exec
UPDATE tasks
SET 
    status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateTaskSpec :exec
UPDATE tasks
SET spec = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: CreateTask :one
INSERT INTO tasks (spec, status, started_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: InsertEvent :one
INSERT INTO events (spec)
VALUES ($1)
RETURNING *;
