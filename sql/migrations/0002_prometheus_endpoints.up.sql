BEGIN;

ALTER TABLE clusters ADD COLUMN prometheus_endpoint TEXT;

COMMIT;
