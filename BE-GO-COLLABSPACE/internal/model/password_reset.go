package model

import (
	"time"

	"github.com/google/uuid"
)

type PasswordReset struct {
	PassID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;column:pass_id"`
	PassUserID    uuid.UUID `gorm:"type:uuid;not null;column:pass_user_id"`
	PassTokenHash string    `gorm:"type:varchar(255);not null;column:pass_token_hash"`
	PassIsUsed    bool      `gorm:"default:false;not null;column:pass_is_used"`
	PassExpireAt  time.Time `gorm:"type:timestamp with time zone;not null;column:pass_expire_at"`
	PassCreatedAt time.Time `gorm:"type:autoCreateTime;column:pass_created_at"`

	User User `gorm:"foreignKey:PassUserID;references:UserID"`
}

func (PasswordReset) TableName() string { return "tbl_password_resets" }
