-- name: GetAutoBackupConfig :one
SELECT * FROM auto_backup_configs
WHERE cluster_id = $1;

-- name: GetAutoDiagnosticsConfig :one
SELECT * FROM auto_diagnostics_configs
WHERE cluster_id = $1;
