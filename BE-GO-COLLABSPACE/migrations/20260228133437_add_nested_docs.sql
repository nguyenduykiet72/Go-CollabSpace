-- +goose Up
SELECT 'up SQL query';
ALTER TABLE tbl_documents ADD CONSTRAINT chk_no_self_parent CHECK (doc_parent_id != doc_id);

CREATE INDEX IF NOT EXISTS idx_docs_workspace_parent
ON tbl_documents (doc_workspace_id, doc_parent_id)
WHERE doc_deleted_at IS NULL;

-- +goose Down
SELECT 'down SQL query';
DROP INDEX IF EXISTS idx_docs_workspace_parent;
ALTER TABLE tbl_documents DROP CONSTRAINT IF EXISTS chk_no_self_parent;
