BEGIN;

CREATE TABLE tasks (
    id          SERIAL PRIMARY KEY,
    attributes  JSONB NOT NULL,
    spec        JSONB NOT NULL,
    status      VARCHAR(255) NOT NULL,
    started_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE events (
    id         SERIAL PRIMARY KEY,
    spec       JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE auto_backup_tasks (
    cluster_id   INTEGER NOT NULL REFERENCES clusters(id),
    next_task_id INTEGER NOT NULL REFERENCES tasks(id),

    UNIQUE (cluster_id)
);

CREATE TABLE auto_diagnostics_tasks (
    cluster_id   INTEGER NOT NULL REFERENCES clusters(id),
    next_task_id INTEGER NOT NULL REFERENCES tasks(id),

    UNIQUE (cluster_id)
);

CREATE TABLE snapshots (
    cluster_id  INTEGER NOT NULL REFERENCES clusters(id),
    snapshot_id BIGINT  NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (cluster_id, snapshot_id)
);

ALTER TABLE auto_backup_configs 
    ADD COLUMN task_id INTEGER NOT NULL REFERENCES tasks(id),
    DROP COLUMN cron_expression,
    DROP COLUMN keep_last;

ALTER TABLE auto_diagnostics_configs 
    ADD COLUMN task_id INTEGER NOT NULL REFERENCES tasks(id),
    DROP COLUMN cron_expression,
    DROP COLUMN retention_duration;

ALTER TABLE organizations
    ADD COLUMN timezone VARCHAR(255) NOT NULL DEFAULT 'UTC';

COMMIT;
