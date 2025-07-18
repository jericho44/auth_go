-- Drop indexes
DROP INDEX IF EXISTS idx_login_attempts_username_attempted_at;
DROP INDEX IF EXISTS idx_login_attempts_attempted_at;
DROP INDEX IF EXISTS idx_login_attempts_username;
DROP INDEX IF EXISTS idx_password_reset_tokens_expires_at;
DROP INDEX IF EXISTS idx_password_reset_tokens_user_id;
DROP INDEX IF EXISTS idx_password_reset_tokens_token;

-- Drop tables
DROP TABLE IF EXISTS login_attempts;
DROP TABLE IF EXISTS password_reset_tokens;