-- name: GetAutoBackupConfig :one
SELECT * FROM auto_backup_configs
WHERE cluster_id = $1;

-- name: GetAutoDiagnosticsConfig :one
SELECT * FROM auto_diagnostics_configs
WHERE cluster_id = $1;

-- name: CreateAutoBackupConfig :exec
INSERT INTO auto_backup_configs (cluster_id, task_id, enabled)
VALUES ($1, $2, $3);

-- name: CreateAutoDiagnosticsConfig :exec
INSERT INTO auto_diagnostics_configs (cluster_id, task_id, enabled)
VALUES ($1, $2, $3);

-- name: UpdateAutoBackupConfig :exec
UPDATE auto_backup_configs
SET enabled = $2
WHERE cluster_id = $1;

-- name: UpdateAutoDiagnosticsConfig :exec
UPDATE auto_diagnostics_configs
SET enabled = $2
WHERE cluster_id = $1;
