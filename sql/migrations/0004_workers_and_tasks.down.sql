BEGIN;

ALTER TABLE auto_backup_configs 
    DROP COLUMN next_task_id;

ALTER TABLE auto_diagnostics_configs 
    DROP COLUMN next_task_id;

ALTER TABLE organizations
    DROP COLUMN timezone;

DROP TABLE auto_backup_tasks;
DROP TABLE auto_diagnostics_tasks;
DROP TABLE tasks;
DROP TABLE events;
DROP TABLE snapshots;

COMMIT;
