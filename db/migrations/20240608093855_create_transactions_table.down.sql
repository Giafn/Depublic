BEGIN;

ALTER TABLE transactions DROP COLUMN auditable;

DROP TABLE IF EXISTS transactions;

COMMIT;
