package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Document struct {
	DocID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;column:doc_id"`
	DocWorkspaceID uuid.UUID  `gorm:"type:uuid;not null;index;column:doc_workspace_id"`
	DocParentID    *uuid.UUID `gorm:"type:uuid;index;column:doc_parent_id"`
	DocAuthorID    uuid.UUID  `gorm:"type:uuid;not null;column:doc_author_id"`
	DocTitle       string     `gorm:"type:varchar(255);not null;default:'Untitled';column:doc_title"`
	DocEmoji       string     `gorm:"type:varchar(10);column:doc_emoji"`
	DocStatus      string     `gorm:"type:varchar(20);default:'private';column:doc_status"`

	DocCreatedAt time.Time      `gorm:"autoCreateTime;column:doc_created_at"`
	DocUpdatedAt time.Time      `gorm:"autoUpdateTime;column:doc_updated_at"`
	DocDeletedAt gorm.DeletedAt `gorm:"index;column:doc_deleted_at"`
}

type DocumentState struct {
	DostDocID     uuid.UUID `gorm:"type:uuid;primaryKey;column:dost_doc_id"`
	DostYjsState  []byte    `gorm:"type:bytea;column:dost_yjs_state"`
	DostPlainText string    `gorm:"type:text;column:dost_plain_text"`
	DostUpdatedAt time.Time `gorm:"autoUpdateTime;column:dost_updated_at"`
}

type DocEmbedding struct {
	EmbID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;column:emb_id"`
	EmbDocID      uuid.UUID `gorm:"type:uuid;not null;index;column:emb_doc_id"`
	EmbContent    string    `gorm:"type:text;not null;column:emb_content"`
	EmbVector     []float32 `gorm:"type:vector(1536);not null;column:emb_vector"`
	EmbTokenCount int       `gorm:"column:emb_token_count"`
}

func (Document) TableName() string      { return "tbl_documents" }
func (DocumentState) TableName() string { return "tbl_document_states" }
func (DocEmbedding) TableName() string  { return "tbl_doc_embeddings" }
