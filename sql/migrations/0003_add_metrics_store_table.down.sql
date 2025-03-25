BEGIN;

DROP TABLE prometheus_endpoints;

ALTER TABLE clusters ADD COLUMN prometheus_endpoint TEXT;

ALTER TABLE clusters DROP COLUMN metrics_store_id;

COMMIT;
