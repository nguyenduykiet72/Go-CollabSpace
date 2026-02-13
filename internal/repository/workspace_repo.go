package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/constant"
	"Go-CollabSpace/internal/model"
)

type workspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(db *gorm.DB) *workspaceRepository {
	return &workspaceRepository{db: db}
}

func (w *workspaceRepository) CreateWorkspace(ctx context.Context, workspace *model.Workspace) error {
	return w.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(workspace).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return apperror.ErrSlugExists
			}
			// var pgErr *pgconn.PgError
			// if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// 	if strings.Contains(pgErr.ConstraintName, "slug") {
			// 		return apperror.ErrSlugExists
			// 	}
			// 	return apperror.ErrSlugExists
			// }
			return err
		}

		var ownerRole model.Role
		if err := tx.Where("role_name = ?", constant.RoleOwner).First(&ownerRole).Error; err != nil {
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

func (w *workspaceRepository) GetWorkspaceByID(ctx context.Context, id uuid.UUID) (*model.Workspace, error) {
	var workspace model.Workspace
	err := w.db.WithContext(ctx).Where("wp_id = ?", id).First(&workspace).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrWorkspaceNotFound
		}
		return nil, err
	}
	return &workspace, nil
}

func (w *workspaceRepository) GetWorkspaceBySlug(ctx context.Context, slug string) (*model.Workspace, error) {
	var workspace model.Workspace
	err := w.db.WithContext(ctx).Where("wp_slug = ?", slug).First(&workspace).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrWorkspaceNotFound
		}
		return nil, err
	}
	return &workspace, nil
}

func (w *workspaceRepository) GetUserWorkspaces(ctx context.Context, userID uuid.UUID) ([]*model.Workspace, error) {
	var workspaces []*model.Workspace
	err := w.db.WithContext(ctx).
		Joins("JOIN tbl_workspace_members ON tbl_workspaces.wp_id = tbl_workspace_members.wpm_workspace_id").
		Where("tbl_workspace_members.wpm_user_id = ?", userID).
		Find(&workspaces).Error

	return workspaces, err
}

func (w *workspaceRepository) IsUserMember(ctx context.Context, workspaceID, userID uuid.UUID) (bool, error) {
	var count int64
	err := w.db.WithContext(ctx).
		Model(&model.WorkspaceMember{}).
		Where("wpm_workspace_id = ? AND wpm_user_id = ?", workspaceID, userID).
		Count(&count).Error

	return count > 0, err
}

func (w *workspaceRepository) GetWorkspaceWithMembers(ctx context.Context, id uuid.UUID) (*model.Workspace, error) {
	var workspace model.Workspace
	err := w.db.WithContext(ctx).
		Preload("Members").
		Preload("Members.User").
		Preload("Members.Role").
		Where("wp_id = ?", id).
		First(&workspace).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrWorkspaceNotFound
		}
		return nil, err
	}
	return &workspace, nil
}

func (w *workspaceRepository) GetRoleByName(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	err := w.db.WithContext(ctx).Where("role_name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrRoleNotFound
		}
		return nil, err
	}
	return &role, nil
}

func (w *workspaceRepository) AddMembers(ctx context.Context, members []model.WorkspaceMember) (int, error) {
	added := 0
	err := w.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, m := range members {
			var count int64
			tx.Model(&model.WorkspaceMember{}).
				Where("wpm_workspace_id = ? AND wpm_user_id = ?", m.WpmWorkspaceID, m.WpmUserID).
				Count(&count)
			if count > 0 {
				continue // Skip existing members
			}

			if err := tx.Create(&m).Error; err != nil {
				return err
			}
			added++
		}
		return nil
	})
	return added, err
}
