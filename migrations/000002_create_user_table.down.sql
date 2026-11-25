DROP INDEX IF EXISTS idx_url_user_id;

DROP INDEX IF EXISTS idx_uq_users_uuid

ALTER TABLE shortener DROP COLUMN IF EXISTS user_id;

DROP TABLE IF EXISTS users; 