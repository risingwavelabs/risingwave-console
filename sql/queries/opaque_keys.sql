-- name: CreateOpaqueKey :one
INSERT INTO opaque_keys (key) VALUES ($1) RETURNING id;

-- name: GetOpaqueKey :one
SELECT key FROM opaque_keys WHERE id = $1;

-- name: DeleteOpaqueKey :exec
DELETE FROM opaque_keys WHERE id = $1;
