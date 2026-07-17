package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;column:user_id"`
	UserEmail    string    `gorm:"type:varchar(255);uniqueIndex;not null;column:user_email"`
	UserFullName string    `gorm:"type:varchar(255);not null;column:user_full_name"`
	UserPassword *string   `gorm:"type:varchar(255);column:user_password"`
	UserAvatar   string    `gorm:"type:text;column:user_avatar"`

	AuthProvider string  `gorm:"type:varchar(50);default:'local';not null;column:auth_provider"`
	SocialID     *string `gorm:"type:varchar(255);uniqueIndex;column:social_id"`

	UserStatus    string         `gorm:"type:varchar(20);default:'active';column:user_status"`
	UserCreatedAt time.Time      `gorm:"type:autoCreateTime;column:user_created_at"`
	UserUpdatedAt time.Time      `gorm:"type:autoUpdateTime;column:user_updated_at"`
	UserDeletedAt gorm.DeletedAt `gorm:"index;column:user_deleted_at"`
}

type Session struct {
	SessID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;column:sess_id"`
	SessUserID       uuid.UUID `gorm:"type:uuid;not null;index;column:sess_user_id"`
	SessRefreshToken string    `gorm:"type:varchar(512);not null;index;column:sess_refresh_token"`
	SessUserAgent    string    `gorm:"type:text;column:sess_user_agent"`
	SessIsBlocked    bool      `gorm:"default:false;column:sess_is_blocked"`
	SessExpireAt     time.Time `gorm:"not null;column:sess_expire_at"`
	SessCreatedAt    time.Time `gorm:"type:autoCreateTime;column:sess_created_at"`

	User User `gorm:"foreignKey:SessUserID;references:UserID"`
}

func (User) TableName() string    { return "tbl_users" }
func (Session) TableName() string { return "tbl_sessions" }
