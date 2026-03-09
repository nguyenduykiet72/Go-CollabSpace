-- +goose Up
SELECT 'up SQL query';
-- Index for full-text search
CREATE INDEX idx_fts_search ON tbl_document_chunks USING GIN (dock_text_search);
-- Index for vector search
CREATE INDEX idx_vector_search ON tbl_document_chunks USING hnsw (dock_embedding vector_cosine_ops);

CREATE TRIGGER tsvector_update_trigger
BEFORE INSERT OR UPDATE
ON tbl_document_chunks
FOR EACH ROW EXECUTE FUNCTION tsvector_update_trigger(dock_text_search,'pg_catalog.english', dock_content);

-- +goose Down
SELECT 'down SQL query';
DROP TRIGGER IF EXISTS tsvector_update_trigger ON tbl_document_chunks;
DROP INDEX IF EXISTS idx_fts_search;
DROP INDEX IF EXISTS idx_vector_search;