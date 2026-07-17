package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateWorkspaceRequest struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
	Slug string `json:"slug" binding:"required,min=3,max=50,alphanum"`
}

type WorkSpaceResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	OwnerID   uuid.UUID `json:"ownerId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type GetWorkspaceParams struct {
	WorkspaceID string `uri:"workspaceId" binding:"required,uuid"`
}

type AddMembersRequest struct {
	UserIDs []uuid.UUID `json:"userIds" binding:"required,min=1,dive,required"`
	Role    string      `json:"role" binding:"required,oneof=Admin Editor Viewer"`
}

type WorkspaceMemberResponse struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"userId"`
	Role     string    `json:"role"`
	JoinedAt string    `json:"joinedAt"`
}
