BEGIN;

CREATE TABLE metrics_stores (
    id              SERIAL PRIMARY KEY,
    name            TEXT NOT NULL,
    spec            JSONB NOT NULL,
    organization_id INTEGER NOT NULL REFERENCES organizations(id),
    default_labels  JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE clusters ADD COLUMN metrics_store_id INTEGER REFERENCES metrics_stores(id);

ALTER TABLE clusters DROP COLUMN prometheus_endpoint;

COMMIT;
