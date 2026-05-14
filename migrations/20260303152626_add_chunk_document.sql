-- +goose Up
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS tbl_document_chunks (
    dock_id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    dock_doc_id UUID NOT NULL,
    dock_content TEXT NOT NULL,
    dock_embedding vector (1536) NOT NULL,
    dock_text_search TSVECTOR,
    dock_created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_chunks_doc FOREIGN KEY (dock_doc_id) REFERENCES tbl_documents (doc_id) ON DELETE CASCADE
);

-- +goose Down
SELECT 'down SQL query';
DROP TABLE IF EXISTS tbl_document_chunks;