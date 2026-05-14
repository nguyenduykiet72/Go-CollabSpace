package dto

import "github.com/google/uuid"

type FileUploadRequest struct {
	ContentType string `json:"content_type" binding:"required"`
	Size        int64  `json:"size" binding:"required"`
}

type FileURLResponse struct {
	UploadURL string    `json:"upload_url"`
	FileID    uuid.UUID `json:"file_id"`
	Key       string    `json:"key"`
	ExpiresIn int       `json:"expires_in"` // Expiration time in seconds
}
