-- name: CreateOrganization :one
INSERT INTO organizations (
    name
) VALUES (
    $1
) ON CONFLICT DO NOTHING RETURNING *;

-- name: GetOrganization :one
SELECT * FROM organizations
WHERE id = $1;

-- name: ListOrganizations :many
SELECT * FROM organizations
ORDER BY name;

-- name: UpdateOrganization :one
UPDATE organizations
SET
    name = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteOrganization :exec
DELETE FROM organizations
WHERE id = $1;

-- name: CreateOrganizationOwner :exec
INSERT INTO organization_owners (
    user_id,
    organization_id
) VALUES (
    $1,
    $2
);
