package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type FileStatus string

const (
	FileStatusPending   FileStatus = "pending"
	FileStatusUploaded  FileStatus = "uploaded"
	FileStatusFailed    FileStatus = "failed"
	FileStatusConfirmed FileStatus = "confirmed"
)

func (f *FileStatus) Scan(value interface{}) error {
	if value == nil {
		*f = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*f = FileStatus(v)
	case string:
		*f = FileStatus(v)
	default:
		return fmt.Errorf("unsupported type for FileStatus: %T", value)
	}
	return nil
}

func (f FileStatus) Value() (driver.Value, error) {
	if string(f) == "" {
		return nil, errors.New("file status is empty")
	}
	return string(f), nil
}

type File struct {
	FileID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;column:file_id"`
	FileUserID    uuid.UUID  `gorm:"type:uuid;not null;index;column:file_user_id"`                        // Uploader's user ID
	FileKey       string     `gorm:"type:varchar(255);not null;column:file_key"`                          // S3 object key
	FileStatus    FileStatus `gorm:"type:file_status_enum;default:'pending';not null;column:file_status"` // pending, uploaded, failed
	FileExpiresAt time.Time  `gorm:"type:timestamp with time zone;not null;column:file_expires_at"`       // Unix timestamp for expiration
	FileCreatedAt time.Time  `gorm:"type:autoCreateTime;column:file_created_at"`

	User User `gorm:"foreignKey:FileUserID;references:UserID"`
}

func (File) TableName() string { return "tbl_files" }
