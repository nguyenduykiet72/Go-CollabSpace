package repository

import (
	"Go-CollabSpace/internal/model"
	"context"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(dbGrm *gorm.DB) *DocumentRepository {
	return &DocumentRepository{db: dbGrm}
}

func (r *DocumentRepository) CreateDoc(ctx context.Context, doc *model.Document) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(doc).Error; err != nil {
			return err
		}

		emptyState := &model.DocumentState{
			DostDocID:     doc.DocID,
			DostPlainText: "",
		}

		return tx.Create(emptyState).Error // Create the initial empty document state
	})
}

func (r *DocumentRepository) GetDocsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.Document, error) {
	var docs []model.Document

	err := r.db.WithContext(ctx).
		Select("doc_id, doc_title, doc_emoji, doc_parent_id, doc_workspace_id").
		Where("doc_workspace_id= ? AND doc_deleted_at IS NULL", workspaceID).
		Find(&docs).Error

	return docs, err
}

func (r *DocumentRepository) GetWorkspaceDocs(ctx context.Context, workspaceID uuid.UUID, userId uuid.UUID) ([]model.Document, error) {
	var docs []model.Document

	err := r.db.WithContext(ctx).
		Where("doc_workspace_id= ? AND doc_author_id = ? AND doc_deleted_at IS NULL", workspaceID, userId).
		Find(&docs).Error

	return docs, err
}

func (r *DocumentRepository) GetDocByID(ctx context.Context, docID uuid.UUID) (*model.Document, error) {
	var doc model.Document
	err := r.db.WithContext(ctx).First(&doc, "doc_id = ?", docID).Error
	return &doc, err
}

func (r *DocumentRepository) UpdateDoc(ctx context.Context, doc *model.Document) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Save(doc).Error
	})
}

func (r *DocumentRepository) DeleteDoc(ctx context.Context, doc *model.Document) error {
	return r.db.WithContext(ctx).Delete(doc).Error
}

func (r *DocumentRepository) AppendYjsUpdate(ctx context.Context, docID uuid.UUID, update []byte) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE tbl_document_states
		 SET dost_yjs_state = COALESCE(dost_yjs_state, '\x'::bytea) || $1
		 WHERE dost_doc_id = $2
		`, update, docID).Error
}

func (r *DocumentRepository) GetYjsState(ctx context.Context, docID uuid.UUID) ([]byte, error) {
	var state model.DocumentState
	err := r.db.WithContext(ctx).
		Select("dost_yjs_state").
		Where("dost_doc_id = ?", docID).
		First(&state).Error
	if err != nil {
		return nil, err
	}

	if state.DostYjsState == nil {
		return []byte{}, nil
	}

	return state.DostYjsState, nil
}

func (r *DocumentRepository) GetDocTreeFlat(ctx context.Context, workspaceID uuid.UUID) ([]model.FlatDocNode, error) {
	query := `
	   WITH RECURSIVE doc_tree as (
	 	 SELECT doc_id, doc_parent_id, doc_title, doc_emoji, doc_status, 0 AS depth
		 FROM tbl_documents 
		 WHERE doc_workspace_id = ?   
		 	AND doc_parent_id IS NULL
		 	AND doc_deleted_at IS NULL

		 UNION ALL

		 SELECT d.doc_id, d.doc_parent_id, d.doc_title, d.doc_emoji, d.doc_status, dt.depth + 1
		 FROM tbl_documents d
		 INNER JOIN doc_tree dt ON d.doc_parent_id = dt.doc_id
		 WHERE d.doc_deleted_at IS NULL
	   )
		SELECT * FROM doc_tree ORDER BY depth, doc_id`

	var nodes []model.FlatDocNode
	if err := r.db.WithContext(ctx).Raw(query, workspaceID).Scan(&nodes).Error; err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r *DocumentRepository) IsDescendant(ctx context.Context, docID, candidateID uuid.UUID) (bool, error) {
	query := `
		WITH RECURSIVE ancestors AS (
			SELECT doc_id
			FROM tbl_documents
			WHERE doc_id = ?
				AND doc_deleted_at IS NULL
			UNION ALL

			SELECT d.doc_id
			FROM tbl_documents d
			INNER JOIN subtree s ON d.doc_parent_id = s.doc_id
			WHERE d.doc_deleted_at IS NULL
		)
		SELECT EXISTS (SELECT 1 FROM subtree WHERE doc_id = ?)
	`

	var exists bool

	if err := r.db.WithContext(ctx).Raw(query, docID, candidateID).Scan(&exists).Error; err != nil {
		return false, err
	}

	return exists, nil
}

func (r *DocumentRepository) UpdateDocParent(ctx context.Context, docID uuid.UUID, newParenID *uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.Document{}).
		Where("doc_id = ? AND doc_deleted_at IS NULL", docID).
		Update("doc_parent_id", newParenID).Error
}

func (r *DocumentRepository) Search(ctx context.Context, keyword string,
	queryEmbedding []float32, workspaceID uuid.UUID, limit int) ([]model.DocumentChunk, error) {
	var chunks []model.DocumentChunk
	vector := pgvector.NewVector(queryEmbedding)

	query := `
		WITH dock_text_search AS (
			SELECT dock_id, ROW_NUMBER() OVER(ORDER BY ts_rank(dock_text_search,plainto_tsquery(?)) DESC) AS rank
			FROM tbl_document_chunks dc
			INNER JOIN tbl_documents d ON dc.dock_doc_id = d.doc_id
			WHERE d.doc_workspace_id = ? AND dock_text_search @@ plainto_tsquery(?) 
			LIMIT 50
		),
		semantic_search AS (
			SELECT dock_id, ROW_NUMBER() OVER(ORDER BY dock_embedding <=> ?) AS rank
			FROM tbl_document_chunks dc 
			INNER JOIN tbl_documents d ON dc.dock_doc_id = d.doc_id
			WHERE d.doc_workspace_id = ?
			LIMIT 50
		)
		SELECT dc.*, COALESCE(1.0 / (60 + ts.rank), 0.0) + COALESCE(1.0 / (60 + ss.rank), 0.0) as rrf_score
		FROM tbl_document_chunks dc
		LEFT JOIN dock_text_search ts ON ts.dock_id = dc.dock_id
		LEFT JOIN semantic_search ss ON ss.dock_id = dc.dock_id
		WHERE ts.dock_id IS NOT NULL OR ss.chunk_id IS NOT NULL
		ORDER BY rrf_score DESC
		LIMIT ?
	`
	err := r.db.WithContext(ctx).Raw(query, keyword, workspaceID, keyword, vector, workspaceID, limit).Scan(&chunks).Error
	return chunks, err
}
