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

-- name: UpdateTask :exec
UPDATE tasks
SET attributes = $2, spec = $3, started_at = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateTaskStartedAt :exec
UPDATE tasks
SET started_at = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: CreateTask :one
INSERT INTO tasks (attributes, spec, status, started_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: InsertEvent :one
INSERT INTO events (spec)
VALUES ($1)
RETURNING *;

-- name: GetTaskByID :one
SELECT * FROM tasks
WHERE id = $1;
