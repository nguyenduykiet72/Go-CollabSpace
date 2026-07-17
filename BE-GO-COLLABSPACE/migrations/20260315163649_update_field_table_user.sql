-- +goose Up
SELECT 'up SQL query';
ALTER TABLE tbl_users
ADD COLUMN auth_provider VARCHAR(50) DEFAULT 'local' NOT NULL,
ADD COLUMN social_id VARCHAR(255) UNIQUE;

ALTER TABLE tbl_users
ALTER COLUMN user_password DROP NOT NULL;

-- +goose Down
SELECT 'down SQL query';
ALTER TABLE tbl_users
ALTER COLUMN user_password SET NOT NULL;
DROP COLUMN auth_provider,
DROP COLUMN social_id;
