-- name: GetTimeZone :one
SELECT timezone
FROM organizations
WHERE id = $1;

-- name: SetTimeZone :exec
UPDATE organizations
SET timezone = $2
WHERE id = $1;
