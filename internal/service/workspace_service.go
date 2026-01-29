package service

import (
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/internal/model"
	"context"

	"github.com/google/uuid"
)

type WorkspaceRepo interface {
	CreateWorkspace(ctx context.Context, workspace *model.Workspace) error
	GetWorkspaceByID(ctx context.Context, id uuid.UUID) (*model.Workspace, error)
	GetWorkspaceBySlug(ctx context.Context, slug string) (*model.Workspace, error)
	GetUserWorkspaces(ctx context.Context, userID uuid.UUID) ([]*model.Workspace, error)
	IsUserMember(ctx context.Context, workspaceID, userID uuid.UUID) (bool, error)
	GetWorkspaceWithMembers(ctx context.Context, id uuid.UUID) (*model.Workspace, error)
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
