package service

import (
	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/internal/model"
	"Go-CollabSpace/internal/worker"
	"context"

	"github.com/google/uuid"
)

type DocumentRepo interface {
	CreateDoc(ctx context.Context, req *model.Document) error
	GetDocsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.Document, error)
	GetDocByID(ctx context.Context, docID uuid.UUID) (*model.Document, error)
	GetDocTreeFlat(ctx context.Context, workspaceID uuid.UUID) ([]model.FlatDocNode, error)
	IsDescendant(ctx context.Context, docID, candidateID uuid.UUID) (bool, error)
	UpdateDocParent(ctx context.Context, docID uuid.UUID, newParentID *uuid.UUID) error
}

type DocumentService struct {
	documentRepo    DocumentRepo
	workspaceRepo   WorkspaceRepo
	taskDistributor worker.TaskDistributor
}

func NewDocumentService(documentRepo DocumentRepo, workspaceRepo WorkspaceRepo, distributor worker.TaskDistributor) *DocumentService {
	return &DocumentService{
		documentRepo:    documentRepo,
		workspaceRepo:   workspaceRepo,
		taskDistributor: distributor,
	}
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

	return &dto.DocumentResponse{
		ID:       doc.DocID,
		Title:    doc.DocTitle,
		Emoji:    doc.DocEmoji,
		ParentID: doc.DocParentID,
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
			ParentID: doc.DocParentID,
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
		ParentID: doc.DocParentID,
	}, nil
}

func (s *DocumentService) GetDocTree(ctx context.Context, workspaceID, userID uuid.UUID) ([]dto.DocTreeItem, error) {
	isMember, err := s.workspaceRepo.IsUserMember(ctx, workspaceID, userID)
	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, apperror.ErrForbidden
	}

	nodes, err := s.documentRepo.GetDocTreeFlat(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	nodeMap := make(map[uuid.UUID]*dto.DocTreeItem, len(nodes))
	for i := range nodes {
		n := nodes[i]
		nodeMap[n.DocID] = &dto.DocTreeItem{
			DocID:    n.DocID,
			ParentID: n.ParentID,
			Title:    n.Title,
			Emoji:    n.Emoji,
			Status:   n.Status,
			Depth:    n.Depth,
			Children: []dto.DocTreeItem{},
		}
	}

	var roots []dto.DocTreeItem
	for i := range nodes {
		n := nodes[i]
		item := nodeMap[n.DocID]
		if n.ParentID == nil {
			roots = append(roots, *item)
		} else {
			parent, ok := nodeMap[*n.ParentID]
			if ok {
				parent.Children = append(parent.Children, *item)
			}
		}
	}

	if roots == nil {
		roots = []dto.DocTreeItem{}
	}

	return roots, nil
}

func (s *DocumentService) MoveDoc(ctx context.Context, docID uuid.UUID, req dto.MoveDocRequest, userID uuid.UUID) error {
	doc, err := s.documentRepo.GetDocByID(ctx, docID)
	if err != nil {
		return err
	}

	isMember, err := s.workspaceRepo.IsUserMember(ctx, doc.DocWorkspaceID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return apperror.ErrForbidden
	}

	if req.NewParentID != nil {
		if *req.NewParentID == docID {
			return apperror.ErrBadRequest
		}

		isDescendant, err := s.documentRepo.IsDescendant(ctx, docID, *req.NewParentID)
		if err != nil {
			return err
		}
		if isDescendant {
			return apperror.ErrCircularReference
		}

		newParent, err := s.documentRepo.GetDocByID(ctx, *req.NewParentID)
		if err != nil {
			return err
		}
		if newParent.DocWorkspaceID != doc.DocWorkspaceID {
			return apperror.ErrForbidden
		}
	}

	return s.documentRepo.UpdateDocParent(ctx, docID, req.NewParentID)
}

func (s *DocumentService) SaveDocSnapshot(ctx context.Context, docID uuid.UUID, req dto.SaveSnapshotRequest) error {
	payload := &worker.PayloadUpdateSearchIndex{
		DocID:     docID,
		PlainText: req.PlainText,
	}

	err := s.taskDistributor.DistributeTaskUpdateSearchIndex(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}
