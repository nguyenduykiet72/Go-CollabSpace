package repository

import (
	"Go-CollabSpace/internal/model"

	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user *model.User) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}
