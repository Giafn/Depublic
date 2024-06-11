BEGIN;

-- tambahkan alamat dan no_hpq
ALTER TABLE users ADD COLUMN alamat TEXT;
ALTER TABLE users ADD COLUMN no_hp VARCHAR(50);

-- rename userId menjadi id
ALTER TABLE users RENAME COLUMN user_id TO id;

COMMIT;