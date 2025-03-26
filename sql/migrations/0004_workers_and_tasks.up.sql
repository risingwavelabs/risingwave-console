BEGIN;

CREATE TABLE workers (
    worker_name    VARCHAR(255),
    last_heartbeat TIMESTAMP NOT NULL,
    UNIQUE (worker_name)
);

CREATE TABLE tasks (
    id          SERIAL PRIMARY KEY,
    worker_name VARCHAR(255),
    spec        JSONB NOT NULL,
    status      VARCHAR(255) NOT NULL,
    timeout     INTERVAL NOT NULL,
    remaining   INT,
    started_at  TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE events (
    id         SERIAL PRIMARY KEY,
    spec       JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE auto_backup_tasks (
    cluster_id INTEGER NOT NULL REFERENCES clusters(id),
    task_id    INTEGER NOT NULL REFERENCES tasks(id),

    UNIQUE (cluster_id)
);

CREATE TABLE auto_diagnostics_tasks (
    cluster_id INTEGER NOT NULL REFERENCES clusters(id),
    task_id    INTEGER NOT NULL REFERENCES tasks(id),

    UNIQUE (cluster_id)
);

COMMIT;
