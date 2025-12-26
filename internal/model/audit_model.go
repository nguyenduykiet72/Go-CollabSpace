package model

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	AudID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;column:aud_id"`
	AudWorkspaceID uuid.UUID `gorm:"type:uuid;index;column:aud_workspace_id"`
	AudActorID     uuid.UUID `gorm:"type:uuid;index;column:aud_actor_id"` // User who performed action

	AudEntity   string `gorm:"type:varchar(50);column:aud_entity"` // E.g., 'document', 'member'
	AudEntityID string `gorm:"type:varchar(50);column:aud_entity_id"`
	AudAction   string `gorm:"type:varchar(50);column:aud_action"` // E.g., 'create', 'delete', 'permission_change'

	AudPayload JSONB `gorm:"type:jsonb;column:aud_payload"`

	AudCreatedAt time.Time `gorm:"autoCreateTime;column:aud_created_at"`
}

type JSONB map[string]interface{}

func (AuditLog) TableName() string { return "tbl_audit_logs" }
