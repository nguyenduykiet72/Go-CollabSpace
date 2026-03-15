-- +goose Up
SELECT 'up SQL query';
ALTER TABLE tbl_password_resets
ADD COLUMN pass_is_used BOOLEAN DEFAULT FALSE NOT NULL;

-- +goose Down
SELECT 'down SQL query';
