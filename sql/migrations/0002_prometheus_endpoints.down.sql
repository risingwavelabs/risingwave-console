BEGIN;

ALTER TABLE clusters DROP COLUMN prometheus_endpoint;

COMMIT;
