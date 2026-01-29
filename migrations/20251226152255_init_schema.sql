-- +goose Up
-- +goose StatementBegin
-- 1. Setup Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "vector";
-- 2. Table: tbl_users
CREATE TABLE IF NOT EXISTS tbl_users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_email VARCHAR(255) NOT NULL,
    user_full_name VARCHAR(255) NOT NULL,
    user_password VARCHAR(255) NOT NULL,
    user_avatar TEXT,
    user_status VARCHAR(20) DEFAULT 'active',
    user_created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    user_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    user_deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX idx_users_email ON tbl_users (user_email);
CREATE INDEX idx_users_deleted_at ON tbl_users (user_deleted_at);
-- 3. Table: tbl_sessions
CREATE TABLE IF NOT EXISTS tbl_sessions (
    sess_id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    sess_user_id UUID NOT NULL,
    sess_refresh_token VARCHAR(512) NOT NULL,
    sess_user_agent TEXT,
    sess_is_blocked BOOLEAN DEFAULT FALSE,
    sess_expire_at TIMESTAMP WITH TIME ZONE NOT NULL,
    sess_created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_sessions_user FOREIGN KEY (sess_user_id) REFERENCES tbl_users (user_id) ON DELETE CASCADE
);
CREATE INDEX idx_sessions_user_id ON tbl_sessions (sess_user_id);
CREATE INDEX idx_sessions_refresh_token ON tbl_sessions (sess_refresh_token);
-- 4. Table: tbl_workspaces
CREATE TABLE IF NOT EXISTS tbl_workspaces (
    wp_id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    wp_owner_id UUID NOT NULL,
    wp_name VARCHAR(100) NOT NULL,
    wp_slug VARCHAR(50) NOT NULL,
    wp_created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    wp_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    wp_deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_workspaces_owner FOREIGN KEY (wp_owner_id) REFERENCES tbl_users (user_id)
);
CREATE UNIQUE INDEX idx_workspaces_slug ON tbl_workspaces (wp_slug);
CREATE INDEX idx_workspaces_deleted_at ON tbl_workspaces (wp_deleted_at);
-- 5. Table: tbl_roles
CREATE TABLE IF NOT EXISTS tbl_roles (
    role_id SERIAL PRIMARY KEY,
    role_name VARCHAR(50) NOT NULL UNIQUE,
    role_description TEXT
);
-- Seed default roles (Optional but recommended)
INSERT INTO tbl_roles (role_name, role_description)
VALUES ('Owner', 'Workspace owner with full access'),
    ('Admin', 'Administrator with management rights'),
    ('Editor', 'Can create and edit content'),
    ('Viewer', 'Read-only access') ON CONFLICT (role_name) DO NOTHING;
-- 6. Table: tbl_workspace_members
CREATE TABLE IF NOT EXISTS tbl_workspace_members (
    wpm_id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    wpm_workspace_id UUID NOT NULL,
    wpm_user_id UUID NOT NULL,
    wpm_role_id INT NOT NULL,
    wpm_joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_members_workspace FOREIGN KEY (wpm_workspace_id) REFERENCES tbl_workspaces (wp_id) ON DELETE CASCADE,
    CONSTRAINT fk_members_user FOREIGN KEY (wpm_user_id) REFERENCES tbl_users (user_id) ON DELETE CASCADE,
    CONSTRAINT fk_members_role FOREIGN KEY (wpm_role_id) REFERENCES tbl_roles (role_id)
);
CREATE UNIQUE INDEX idx_wp_user ON tbl_workspace_members (wpm_workspace_id, wpm_user_id);
-- 7. Table: tbl_documents
CREATE TABLE IF NOT EXISTS tbl_documents (
    doc_id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    doc_workspace_id UUID NOT NULL,
    doc_parent_id UUID,
    doc_author_id UUID NOT NULL,
    doc_title VARCHAR(255) NOT NULL DEFAULT 'Untitled',
    doc_emoji VARCHAR(10),
    doc_status VARCHAR(20) DEFAULT 'private',
    doc_created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    doc_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    doc_deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_docs_workspace FOREIGN KEY (doc_workspace_id) REFERENCES tbl_workspaces (wp_id) ON DELETE CASCADE,
    CONSTRAINT fk_docs_author FOREIGN KEY (doc_author_id) REFERENCES tbl_users (user_id),
    CONSTRAINT fk_docs_parent FOREIGN KEY (doc_parent_id) REFERENCES tbl_documents (doc_id) ON DELETE
    SET NULL
);
CREATE INDEX idx_docs_workspace_id ON tbl_documents (doc_workspace_id);
CREATE INDEX idx_docs_parent_id ON tbl_documents (doc_parent_id);
CREATE INDEX idx_docs_deleted_at ON tbl_documents (doc_deleted_at);
-- 8. Table: tbl_document_states (Yjs Binary Storage)
CREATE TABLE IF NOT EXISTS tbl_document_states (
    dost_doc_id UUID PRIMARY KEY,
    dost_yjs_state BYTEA,
    -- Dữ liệu binary của Yjs
    dost_plain_text TEXT,
    -- Dữ liệu text để preview/fallback
    dost_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_states_doc FOREIGN KEY (dost_doc_id) REFERENCES tbl_documents (doc_id) ON DELETE CASCADE
);
-- 9. Table: tbl_doc_embeddings (Vector Search)
CREATE TABLE IF NOT EXISTS tbl_doc_embeddings (
    emb_id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    emb_doc_id UUID NOT NULL,
    emb_content TEXT NOT NULL,
    emb_vector vector (1536) NOT NULL,
    -- OpenAI model dimension
    emb_token_count INT,
    CONSTRAINT fk_embeddings_doc FOREIGN KEY (emb_doc_id) REFERENCES tbl_documents (doc_id) ON DELETE CASCADE
);
CREATE INDEX idx_embeddings_doc_id ON tbl_doc_embeddings (emb_doc_id);
CREATE INDEX idx_embeddings_vector ON tbl_doc_embeddings USING hnsw (emb_vector vector_cosine_ops);
-- 10. Table: tbl_audit_logs
CREATE TABLE IF NOT EXISTS tbl_audit_logs (
    aud_id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    aud_workspace_id UUID,
    aud_actor_id UUID,
    aud_entity VARCHAR(50),
    aud_entity_id VARCHAR(50),
    aud_action VARCHAR(50),
    aud_payload JSONB,
    aud_created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_audit_workspace ON tbl_audit_logs (aud_workspace_id);
CREATE INDEX idx_audit_actor ON tbl_audit_logs (aud_actor_id);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tbl_audit_logs;
DROP TABLE IF EXISTS tbl_doc_embeddings;
DROP TABLE IF EXISTS tbl_document_states;
DROP TABLE IF EXISTS tbl_documents;
DROP TABLE IF EXISTS tbl_workspace_members;
DROP TABLE IF EXISTS tbl_roles;
DROP TABLE IF EXISTS tbl_workspaces;
DROP TABLE IF EXISTS tbl_sessions;
DROP TABLE IF EXISTS tbl_users;
-- Keep extensions unless you really want to remove them (usually kept for other tables)
-- DROP EXTENSION IF EXISTS "vector";
-- DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd