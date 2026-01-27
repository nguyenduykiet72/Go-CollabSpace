package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(dbGrm *gorm.DB) *UserRepository {
	return &UserRepository{db: dbGrm}
}

func (u *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	return u.db.WithContext(ctx).Create(user).Error
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).Where("user_email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).First(&user, "user_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) CreateSession(ctx context.Context, session *model.Session) error {
	return u.db.WithContext(ctx).Create(session).Error
}

func (u *UserRepository) RevokeSession(ctx context.Context, id uuid.UUID) error {
	return u.db.WithContext(ctx).Model(&model.Session{}).Where("sess_id = ?", id).Update("sess_is_blocked", true).Error
}
