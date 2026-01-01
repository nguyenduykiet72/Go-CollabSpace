package repository

import (
	"Go-CollabSpace/internal/model"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)

	CreateSession(ctx context.Context, session *model.Session) error
	RevokeSession(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}

func (u userRepository) CreateUser(ctx context.Context, user *model.User) error {
	return u.db.WithContext(ctx).Create(user).Error
}

func (u userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).Where("user_email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).First(&user, "user_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u userRepository) CreateSession(ctx context.Context, session *model.Session) error {
	return u.db.WithContext(ctx).Create(session).Error
}

func (u userRepository) RevokeSession(ctx context.Context, id uuid.UUID) error {
	return u.db.WithContext(ctx).Model(&model.Session{}).Where("sess_id = ?", id).Update("sess_is_blocked", true).Error
}
