-- name: GetOrgSettings :one
SELECT * FROM org_settings
WHERE org_id = $1;

-- name: CreateOrgSettings :exec
INSERT INTO org_settings (org_id, timezone)
VALUES ($1, $2);
