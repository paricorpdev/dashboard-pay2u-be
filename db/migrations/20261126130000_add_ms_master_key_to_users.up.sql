ALTER TABLE users
    ADD COLUMN IF NOT EXISTS ms_master_key VARCHAR(128);

UPDATE users u
SET ms_master_key = 'dba_' || replace(gen_random_uuid()::text, '-', '')
WHERE u.ms_master_key IS NULL
  AND EXISTS (
      SELECT 1
      FROM user_roles ur
      JOIN roles r ON r.id = ur.role_id
      WHERE ur.user_id = u.id
        AND r.code IN ('ROLE_ADMIN', 'ROLE_SUPER_ADMIN')
  );

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_ms_master_key
    ON users (ms_master_key)
    WHERE ms_master_key IS NOT NULL;
