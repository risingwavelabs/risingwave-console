BEGIN;

-- migrate tasks

INSERT INTO anchor.tasks (id, attributes, spec, status, started_at, created_at, updated_at) SELECT
    id,
    attributes,
    spec,
    status,
    started_at,
    created_at,
    updated_at
FROM tasks;

INSERT INTO anchor.events (id, spec, created_at) SELECT
    id,
    spec,
    created_at
FROM events;

DROP TABLE auto_backup_tasks;
DROP TABLE auto_diagnostics_tasks;

ALTER TABLE auto_backup_configs
    DROP CONSTRAINT auto_backup_configs_task_id_fkey;

ALTER TABLE auto_backup_configs
    ADD CONSTRAINT auto_backup_configs_task_id_fkey FOREIGN KEY (task_id) REFERENCES anchor.tasks(id);

ALTER TABLE auto_diagnostics_configs
    DROP CONSTRAINT auto_diagnostics_configs_task_id_fkey;

ALTER TABLE auto_diagnostics_configs
    ADD CONSTRAINT auto_diagnostics_configs_task_id_fkey FOREIGN KEY (task_id) REFERENCES anchor.tasks(id);

DROP TABLE events;
DROP TABLE tasks;

-- migrate organizations and users

INSERT INTO anchor.orgs (id, name, tz, created_at, updated_at) SELECT
    id,
    name,
    timezone,
    created_at,
    updated_at
FROM organizations;

INSERT INTO anchor.users (id, name, password_hash, password_salt, created_at, updated_at) SELECT
    id,
    name,
    password_hash,
    password_salt,
    created_at,
    updated_at
FROM users;

INSERT INTO anchor.org_owners (org_id, user_id, created_at) SELECT
    organization_id,
    user_id,
    created_at
FROM organization_owners;

INSERT INTO anchor.org_users (org_id, user_id, created_at) SELECT
    organization_id,
    id,
    created_at
FROM users;

-- migrate org fk

ALTER TABLE clusters
    DROP CONSTRAINT clusters_organization_id_fkey;

ALTER TABLE clusters
    RENAME COLUMN organization_id TO org_id;

ALTER TABLE clusters
    ADD CONSTRAINT clusters_orgs_id_fkey FOREIGN KEY (org_id) REFERENCES anchor.orgs(id);

ALTER TABLE database_connections
    DROP CONSTRAINT database_connections_organization_id_fkey;

ALTER TABLE database_connections
    RENAME COLUMN organization_id TO org_id;

ALTER TABLE database_connections
    ADD CONSTRAINT database_connections_orgs_id_fkey FOREIGN KEY (org_id) REFERENCES anchor.orgs(id);

ALTER TABLE metrics_stores 
    DROP CONSTRAINT metrics_stores_organization_id_fkey;

ALTER TABLE metrics_stores
    RENAME COLUMN organization_id TO org_id;

ALTER TABLE metrics_stores
    ADD CONSTRAINT metrics_stores_orgs_id_fkey FOREIGN KEY (org_id) REFERENCES anchor.orgs(id);

ALTER TABLE org_settings
    DROP CONSTRAINT org_settings_organization_id_fkey;

ALTER TABLE org_settings
    RENAME COLUMN organization_id TO org_id;

ALTER TABLE org_settings
    ADD CONSTRAINT org_settings_orgs_id_fkey FOREIGN KEY (org_id) REFERENCES anchor.orgs(id);

COMMIT;
