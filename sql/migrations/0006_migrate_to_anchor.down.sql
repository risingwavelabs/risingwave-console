BEGIN;

-- revert org fk changes
ALTER TABLE org_settings
    DROP CONSTRAINT org_settings_orgs_id_fkey;

ALTER TABLE org_settings
    RENAME COLUMN org_id TO organization_id;

ALTER TABLE org_settings
    ADD CONSTRAINT org_settings_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES organizations(id);

ALTER TABLE metrics_stores
    DROP CONSTRAINT metrics_stores_orgs_id_fkey;

ALTER TABLE metrics_stores
    RENAME COLUMN org_id TO organization_id;

ALTER TABLE metrics_stores
    ADD CONSTRAINT metrics_stores_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES organizations(id);

ALTER TABLE database_connections
    DROP CONSTRAINT database_connections_orgs_id_fkey;

ALTER TABLE database_connections
    RENAME COLUMN org_id TO organization_id;

ALTER TABLE database_connections
    ADD CONSTRAINT database_connections_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES organizations(id);

ALTER TABLE clusters
    DROP CONSTRAINT clusters_orgs_id_fkey;

ALTER TABLE clusters
    RENAME COLUMN org_id TO organization_id;

ALTER TABLE clusters
    ADD CONSTRAINT clusters_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES organizations(id);

-- recreate original tables
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY,
    attributes JSONB NOT NULL,
    spec JSONB NOT NULL,
    status VARCHAR(255) NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY,
    spec JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- restore data from anchor schema
INSERT INTO tasks (id, attributes, spec, status, started_at, created_at, updated_at) 
SELECT id, attributes, spec, status, started_at, created_at, updated_at
FROM anchor.tasks;

INSERT INTO events (id, spec, created_at) 
SELECT id, spec, created_at
FROM anchor.events;

-- restore foreign keys for tasks
ALTER TABLE auto_backup_configs
    DROP CONSTRAINT auto_backup_configs_task_id_fkey;

ALTER TABLE auto_backup_configs
    ADD CONSTRAINT auto_backup_configs_task_id_fkey FOREIGN KEY (task_id) REFERENCES tasks(id);

ALTER TABLE auto_diagnostics_configs
    DROP CONSTRAINT auto_diagnostics_configs_task_id_fkey;

ALTER TABLE auto_diagnostics_configs
    ADD CONSTRAINT auto_diagnostics_configs_task_id_fkey FOREIGN KEY (task_id) REFERENCES tasks(id);

-- recreate organizations and users tables
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    timezone VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    password_salt VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS organization_owners (
    organization_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (organization_id, user_id),
    FOREIGN KEY (organization_id) REFERENCES organizations(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- restore data from anchor schema
INSERT INTO organizations (id, name, timezone, created_at, updated_at)
SELECT id, name, tz, created_at, updated_at
FROM anchor.orgs;

INSERT INTO users (id, name, password_hash, password_salt, created_at, updated_at)
SELECT id, name, password_hash, password_salt, created_at, updated_at
FROM anchor.users;

INSERT INTO organization_owners (organization_id, user_id, created_at)
SELECT org_id, user_id, created_at
FROM anchor.org_owners;

COMMIT;
