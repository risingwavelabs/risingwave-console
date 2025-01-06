-- name: UpsertRefreshToken :exec
INSERT INTO refresh_tokens (user_id, token)
VALUES ($1, $2) ON CONFLICT (user_id, token) DO UPDATE SET token = EXCLUDED.token, updated_at = CURRENT_TIMESTAMP;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE user_id = $1 AND token = $2;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE user_id = $1 AND token = $2;
