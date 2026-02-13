package service

import (
	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/internal/model"
	"context"

	"github.com/google/uuid"
)

type DocumentRepo interface {
	CreateDoc(ctx context.Context, req *model.Document) error
	GetDocsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.Document, error)
	GetDocByID(ctx context.Context, docID uuid.UUID) (*model.Document, error)
}

type DocumentService struct {
	documentRepo  DocumentRepo
	workspaceRepo WorkspaceRepo
}

func NewDocumentService(documentRepo DocumentRepo, workspaceRepo WorkspaceRepo) *DocumentService {
	return &DocumentService{documentRepo: documentRepo, workspaceRepo: workspaceRepo}
}

func (s *DocumentService) CreateDoc(ctx context.Context, req dto.CreateDocRequest, userID uuid.UUID) (*dto.DocumentResponse, error) {
	isMember, _ := s.workspaceRepo.IsUserMember(ctx, req.WorkspaceID, userID)
	if !isMember {
		return nil, apperror.ErrUnauthorized
	}

	doc := &model.Document{
		DocWorkspaceID: req.WorkspaceID,
		DocAuthorID:    userID,
		DocParentID:    req.ParentID,
		DocTitle:       req.Title,
		DocEmoji:       req.Emoji,
		DocStatus:      "active",
	}

	if err := s.documentRepo.CreateDoc(ctx, doc); err != nil {
		return nil, err
	}

	var parentID uuid.UUID
	if doc.DocParentID != nil {
		parentID = *doc.DocParentID
	}

	return &dto.DocumentResponse{
		ID:       doc.DocID,
		Title:    doc.DocTitle,
		Emoji:    doc.DocEmoji,
		ParentID: parentID,
	}, nil
}

func (s *DocumentService) GetWorkspaceDocs(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID) ([]dto.DocumentResponse, error) {
	isMember, _ := s.workspaceRepo.IsUserMember(ctx, workspaceID, userID)
	if !isMember {
		return nil, apperror.ErrUnauthorized
	}

	docs, err := s.documentRepo.GetDocsByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	var response []dto.DocumentResponse
	for _, doc := range docs {
		response = append(response, dto.DocumentResponse{
			ID:       doc.DocID,
			Title:    doc.DocTitle,
			Emoji:    doc.DocEmoji,
			ParentID: *doc.DocParentID,
		})
	}

	return response, nil
}

func (s *DocumentService) GetDocDetail(ctx context.Context, docID uuid.UUID, userID uuid.UUID) (*dto.DocumentResponse, error) {
	doc, err := s.documentRepo.GetDocByID(ctx, docID)
	if err != nil {
		return nil, err
	}

	isMember, _ := s.workspaceRepo.IsUserMember(ctx, doc.DocWorkspaceID, userID)
	if !isMember {
		return nil, apperror.ErrUnauthorized
	}

	return &dto.DocumentResponse{
		ID:       doc.DocID,
		Title:    doc.DocTitle,
		Emoji:    doc.DocEmoji,
		ParentID: *doc.DocParentID,
	}, nil
}
