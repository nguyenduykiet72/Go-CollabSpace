package service

import (
	"Go-CollabSpace/internal/model"
	"Go-CollabSpace/internal/repository"
)

type IUserService interface {
	Register(name string, email string) error
}

type UserService struct {
	repo repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) IUserService {
	return &UserService{repo: repo}
}

func (u UserService) Register(name string, email string) error {
	user := &model.User{
		Name:  name,
		Email: email,
	}
	return u.repo.CreateUser(user)
}
