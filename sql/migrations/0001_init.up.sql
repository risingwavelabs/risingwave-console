BEGIN;

CREATE TABLE IF NOT EXISTS organizations (
    id            SERIAL,
    name          TEXT        NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    UNIQUE (name),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS users (
    id              SERIAL,
    name            TEXT        NOT NULL,
    password_hash   TEXT        NOT NULL,
    password_salt   TEXT        NOT NULL,
    organization_id INTEGER     NOT NULL REFERENCES organizations(id),
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    UNIQUE (name),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS organization_owners (
    user_id         INTEGER     NOT NULL REFERENCES users(id),
    organization_id INTEGER     NOT NULL REFERENCES organizations(id),
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    PRIMARY KEY (user_id, organization_id)
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id              SERIAL,
    user_id         INTEGER     NOT NULL REFERENCES users(id),
    token           TEXT        NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    PRIMARY KEY (id),
    UNIQUE (user_id, token)
);

CREATE TABLE IF NOT EXISTS clusters (
    id              SERIAL,
    organization_id INTEGER     NOT NULL REFERENCES organizations(id),

    name            TEXT        NOT NULL,
    host            TEXT        NOT NULL,
    sql_port        INTEGER     NOT NULL,
    meta_port       INTEGER     NOT NULL,
    
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    UNIQUE (organization_id, name),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS database_connections (
    id              SERIAL,
    organization_id INTEGER     NOT NULL REFERENCES organizations(id),

    name            TEXT        NOT NULL,
    cluster_id      INTEGER     NOT NULL REFERENCES clusters(id),
    username        TEXT,
    password        TEXT,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    UNIQUE (name),
    PRIMARY KEY (id)
);

COMMIT;
