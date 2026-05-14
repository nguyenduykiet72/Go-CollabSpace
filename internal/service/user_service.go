package service

import (
	"context"
	"time"

	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/internal/model"
)

type UserRepo interface {
	GetAllUsers(ctx context.Context, req dto.PaginationReq) ([]*model.User, error)
}

type UserService struct {
	userRepo   UserRepo
	transactor Transactor
}

func NewUserService(userRepo UserRepo, transactor Transactor) *UserService {
	return &UserService{userRepo: userRepo, transactor: transactor}
}

func (u *UserService) GetAllUsers(ctx context.Context, req dto.PaginationReq) ([]*dto.UserResponse, error) {
	users, err := u.userRepo.GetAllUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	var userResponses []*dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, &dto.UserResponse{
			ID:        user.UserID,
			Email:     user.UserEmail,
			FullName:  user.UserFullName,
			Avatar:    user.UserAvatar,
			CreatedAt: user.UserCreatedAt.Format(time.RFC3339),
		})
	}
	return userResponses, nil
}
