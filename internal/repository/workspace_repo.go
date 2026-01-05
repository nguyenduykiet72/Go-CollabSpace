package repository

import (
	"Go-CollabSpace/internal/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type IWorkspaceRepository interface {
	CreateWorkSpace(ctx context.Context, workspace *model.Workspace) error
	GetRoleByName(ctx context.Context, name string) (*model.Role, error)
}

type workspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(db *gorm.DB) IWorkspaceRepository {
	return &workspaceRepository{db: db}
}

func (w workspaceRepository) CreateWorkSpace(ctx context.Context, workspace *model.Workspace) error {
	return w.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(workspace).Error; err != nil {
			return err
		}

		var ownerRole model.Role
		if err := tx.Where("role_name = ?", "Owner").First(&ownerRole).Error; err != nil {
			return errors.New("system error: Owner role not found")
		}

		member := model.WorkspaceMember{
			WpmWorkspaceID: workspace.WpID,
			WpmUserID:      workspace.WpOwnerID,
			WpmRoleID:      ownerRole.RoleID,
		}

		if err := tx.Create(&member).Error; err != nil {
			return err
		}
		return nil
	})
}

func (w workspaceRepository) GetRoleByName(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	if err := w.db.WithContext(ctx).Where("role_name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
