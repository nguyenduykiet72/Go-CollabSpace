-- +goose Up
SELECT 'up SQL query';

CREATE TYPE file_status_enum AS ENUM ('pending','uploaded','failed','confirmed');

CREATE TABLE IF NOT EXISTS tbl_files (
    file_id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    file_user_id UUID NOT NULL,
    file_key VARCHAR(255) NOT NULL,
    file_status file_status_enum DEFAULT 'pending' NOT NULL,
    file_expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    file_created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_files_status_expires ON tbl_files (file_status, file_expires_at);

-- +goose Down
SELECT 'down SQL query';

DROP TABLE IF EXISTS tbl_files;
DROP TYPE IF EXISTS file_status_enum;