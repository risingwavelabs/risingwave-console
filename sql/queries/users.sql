-- name: CreateUser :one
INSERT INTO users (
    name,
    password_hash,
    password_salt,
    organization_id
) VALUES (
    $1, $2, $3, $4
) ON CONFLICT (name) DO UPDATE 
    SET updated_at = CURRENT_TIMESTAMP, 
        name = EXCLUDED.name, 
        password_hash = EXCLUDED.password_hash,
        password_salt = EXCLUDED.password_salt
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByName :one
SELECT * FROM users
WHERE name = $1;

-- name: DeleteUserByName :exec
DELETE FROM users
WHERE name = $1;
