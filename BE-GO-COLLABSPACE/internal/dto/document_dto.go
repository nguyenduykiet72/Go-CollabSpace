package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateDocRequest struct {
	WorkspaceID uuid.UUID  `json:"workspaceId" binding:"required"`
	ParentID    *uuid.UUID `json:"parentId,omitempty"`
	Title       string     `json:"title" binding:"required,min=1,max=200"`
	Emoji       string     `json:"emoji"`
}

type DocumentResponse struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Emoji     string     `json:"emoji"`
	ParentID  *uuid.UUID `json:"parentId"` // pointer: nil = root document
	AuthorID  uuid.UUID  `json:"authorId"`
	CreatedAt time.Time  `json:"createdAt"`
}

type DocTreeItem struct {
	DocID    uuid.UUID     `json:"docId"`
	Title    string        `json:"title"`
	ParentID *uuid.UUID    `json:"parentId"`
	Children []DocTreeItem `json:"children"` // Recursive definition for child documents
	Emoji    string        `json:"emoji"`
	Status   string        `json:"status"`
	Depth    int           `json:"depth"` // Optional: to track the depth in the tree
}

type UpdateDocRequest struct {
	Title string `json:"title"`
	Emoji string `json:"emoji"`
}

type MoveDocRequest struct {
	NewParentID *uuid.UUID `json:"newParentId"` // Nullable for moving to root
}

type SaveSnapshotRequest struct {
	PlainText string `json:"plainText" binding:"required"`
}
