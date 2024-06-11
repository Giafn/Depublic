BEGIN;

-- hapus alamat dan no_hp
ALTER TABLE users DROP COLUMN alamat;
ALTER TABLE users DROP COLUMN no_hp;

-- rename id menjadi userId
ALTER TABLE users RENAME COLUMN id TO user_id;

COMMIT;