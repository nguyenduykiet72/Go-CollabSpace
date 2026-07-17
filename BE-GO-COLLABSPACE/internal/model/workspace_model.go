package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Workspace struct {
	WpID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;column:wp_id"`
	WpOwnerID uuid.UUID `gorm:"type:uuid;not null;column:wp_owner_id"`
	WpName    string    `gorm:"type:varchar(100);not null;column:wp_name"`
	WpSlug    string    `gorm:"type:varchar(50);uniqueIndex;not null;column:wp_slug"`

	WpCreatedAt time.Time      `gorm:"autoCreateTime;column:wp_created_at"`
	WpUpdatedAt time.Time      `gorm:"autoUpdateTime;column:wp_updated_at"`
	WpDeletedAt gorm.DeletedAt `gorm:"index;column:wp_deleted_at"`

	Owner  User              `gorm:"foreignKey:WpOwnerID;references:UserID"`
	Member []WorkspaceMember `gorm:"foreignKey:WpmWorkspaceID;references:WpID"`
}

type Role struct {
	RoleID          uint   `gorm:"primaryKey;autoIncrement;column:role_id"`
	RoleName        string `gorm:"type:varchar(50);unique;not null;column:role_name"`
	RoleDescription string `gorm:"type:text;column:role_description"`
}

type WorkspaceMember struct {
	WpmID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;column:wpm_id"`
	WpmWorkspaceID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_wp_user;column:wpm_workspace_id"`
	WpmUserID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_wp_user;column:wpm_user_id"`
	WpmRoleID      uint      `gorm:"not null;column:wpm_role_id"`

	WpmJoinedAt time.Time `gorm:"autoCreateTime;column:wpm_joined_at"`

	Workspace Workspace `gorm:"foreignKey:WpmWorkspaceID;references:WpID"`
	User      User      `gorm:"foreignKey:WpmUserID;references:UserID"`
	Role      Role      `gorm:"foreignKey:WpmRoleID;references:RoleID"`
}

func (Workspace) TableName() string       { return "tbl_workspaces" }
func (Role) TableName() string            { return "tbl_roles" }
func (WorkspaceMember) TableName() string { return "tbl_workspace_members" }
