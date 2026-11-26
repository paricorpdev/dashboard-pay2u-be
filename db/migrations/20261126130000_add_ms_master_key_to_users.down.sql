DROP INDEX IF EXISTS idx_users_ms_master_key;

ALTER TABLE users
    DROP COLUMN IF EXISTS ms_master_key;
