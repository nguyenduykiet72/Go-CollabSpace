package repository

import (
	"Go-CollabSpace/internal/model"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type documentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(dbGrm *gorm.DB) *documentRepository {
	return &documentRepository{db: dbGrm}
}

func (r *documentRepository) CreateDoc(ctx context.Context, doc *model.Document) error {
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

func (r *documentRepository) GetDocsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.Document, error) {
	var docs []model.Document

	err := r.db.WithContext(ctx).
		Select("doc_id, doc_title, doc_emoji, doc_parent_id, doc_workspace_id").
		Where("doc_workspace_id= ? AND doc_deleted_at IS NULL", workspaceID).
		Find(&docs).Error

	return docs, err
}

func (r *documentRepository) GetWorkspaceDocs(ctx context.Context, workspaceID uuid.UUID, userId uuid.UUID) ([]model.Document, error) {
	var docs []model.Document

	err := r.db.WithContext(ctx).
		Where("doc_workspace_id= ? AND doc_author_id = ? AND doc_deleted_at IS NULL", workspaceID, userId).
		Find(&docs).Error

	return docs, err
}

func (r *documentRepository) GetDocByID(ctx context.Context, docID uuid.UUID) (*model.Document, error) {
	var doc model.Document
	err := r.db.WithContext(ctx).First(&doc, "doc_id = ?", docID).Error
	return &doc, err
}
