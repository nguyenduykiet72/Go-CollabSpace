-- +goose Up
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS tbl_password_resets (
    pass_id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    pass_user_id UUID NOT NULL,
    pass_token_hash VARCHAR(255) NOT NULL,
    pass_expire_at TIMESTAMP WITH TIME ZONE NOT NULL,
    pass_created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_password_resets_user FOREIGN KEY (pass_user_id) REFERENCES tbl_users (user_id) ON DELETE CASCADE
);

-- +goose Down
SELECT 'down SQL query';
DROP TABLE IF EXISTS tbl_password_resets;