package service

import (
	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/internal/model"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type WorkspaceRepo interface {
	CreateWorkspace(ctx context.Context, workspace *model.Workspace) error
	GetWorkspaceByID(ctx context.Context, id uuid.UUID) (*model.Workspace, error)
	GetWorkspaceBySlug(ctx context.Context, slug string) (*model.Workspace, error)
	GetUserWorkspaces(ctx context.Context, userID uuid.UUID) ([]*model.Workspace, error)
	IsUserMember(ctx context.Context, workspaceID, userID uuid.UUID) (bool, error)
	GetWorkspaceWithMembers(ctx context.Context, id uuid.UUID) (*model.Workspace, error)
	GetRoleByName(ctx context.Context, name string) (*model.Role, error)
	AddMembers(ctx context.Context, members []model.WorkspaceMember) (int, error)
}

type WorkspaceService struct {
	workspaceRepo WorkspaceRepo
}

func NewWorkspaceService(workspaceRepo WorkspaceRepo) *WorkspaceService {
	return &WorkspaceService{workspaceRepo: workspaceRepo}
}

func (s *WorkspaceService) CreateWorkspace(ctx context.Context, req dto.CreateWorkspaceRequest, userID uuid.UUID) (*dto.WorkSpaceResponse, error) {
	newWorkspace := &model.Workspace{
		WpName:    req.Name,
		WpSlug:    req.Slug,
		WpOwnerID: userID,
	}

	if err := s.workspaceRepo.CreateWorkspace(ctx, newWorkspace); err != nil {
		return nil, err
	}

	return &dto.WorkSpaceResponse{
		ID:        newWorkspace.WpID,
		Name:      newWorkspace.WpName,
		Slug:      newWorkspace.WpSlug,
		OwnerID:   newWorkspace.WpOwnerID,
		CreatedAt: newWorkspace.WpCreatedAt,
	}, nil
}

func (s *WorkspaceService) GetWorkspaceByID(ctx context.Context, id uuid.UUID) (*dto.WorkSpaceResponse, error) {
	workspace, err := s.workspaceRepo.GetWorkspaceByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.WorkSpaceResponse{
		ID:        workspace.WpID,
		Name:      workspace.WpName,
		Slug:      workspace.WpSlug,
		OwnerID:   workspace.WpOwnerID,
		CreatedAt: workspace.WpCreatedAt,
		UpdatedAt: workspace.WpUpdatedAt,
	}, nil
}

func (s *WorkspaceService) AddMembers(ctx context.Context, workspaceID uuid.UUID, req dto.AddMembersRequest, callerID uuid.UUID) (int, error) {
	// Check caller is a member of the workspace
	isMember, _ := s.workspaceRepo.IsUserMember(ctx, workspaceID, callerID)
	if !isMember {
		return 0, apperror.ErrUnauthorized
	}

	// Get the role
	role, err := s.workspaceRepo.GetRoleByName(ctx, req.Role)
	if err != nil {
		return 0, err
	}

	// Build member models
	members := make([]model.WorkspaceMember, 0, len(req.UserIDs))
	for _, userID := range req.UserIDs {
		members = append(members, model.WorkspaceMember{
			WpmWorkspaceID: workspaceID,
			WpmUserID:      userID,
			WpmRoleID:      role.RoleID,
		})
	}

	added, err := s.workspaceRepo.AddMembers(ctx, members)
	if err != nil {
		return 0, fmt.Errorf("failed to add members: %w", err)
	}

	return added, nil
}
