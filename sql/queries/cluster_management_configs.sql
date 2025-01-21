-- name: UpsertAutoBackupConfig :exec
INSERT INTO auto_backup_configs (cluster_id, enabled, cron_expression, keep_last)
VALUES ($1, $2, $3, $4)
ON CONFLICT (cluster_id) DO UPDATE SET enabled = $2, cron_expression = $3, keep_last = $4;

-- name: UpsertAutoDiagnosticsConfig :exec
INSERT INTO auto_diagnostics_configs (cluster_id, enabled, cron_expression, retention_duration)
VALUES ($1, $2, $3, $4)
ON CONFLICT (cluster_id) DO UPDATE SET enabled = $2, cron_expression = $3, retention_duration = $4;

-- name: GetAutoBackupConfig :one
SELECT * FROM auto_backup_configs
WHERE cluster_id = $1;

-- name: GetAutoDiagnosticsConfig :one
SELECT * FROM auto_diagnostics_configs
WHERE cluster_id = $1;
