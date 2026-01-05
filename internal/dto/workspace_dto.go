package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateWorkspaceRequest struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
	Slug string `json:"name" binding:"required,min=3,max=50,alphanum"`
}

type WorkSpaceResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	OwnerID   uuid.UUID `json:"ownerId"`
	CreatedAt time.Time `json:"createdAt"`
}
