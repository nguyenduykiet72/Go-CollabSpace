package service

import (
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/internal/model"
	"Go-CollabSpace/internal/repository"
	"context"

	"github.com/google/uuid"
)

type IWorkspaceService interface {
	Create(ctx context.Context, req dto.CreateWorkspaceRequest, userID uuid.UUID) (*dto.WorkSpaceResponse, error)
}

type WorkspaceService struct {
	repository.IWorkspaceRepository
}

func NewWorkspaceService(workspaceRepo repository.IWorkspaceRepository) IWorkspaceService {
	return &WorkspaceService{IWorkspaceRepository: workspaceRepo}
}

func (s *WorkspaceService) Create(ctx context.Context, req dto.CreateWorkspaceRequest, userID uuid.UUID) (*dto.WorkSpaceResponse, error) {
	newWorkspace := &model.Workspace{
		WpName:    req.Name,
		WpSlug:    req.Slug,
		WpOwnerID: userID,
	}

	if err := s.CreateWorkSpace(ctx, newWorkspace); err != nil {
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
