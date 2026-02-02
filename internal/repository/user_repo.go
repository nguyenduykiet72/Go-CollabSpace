package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/internal/model"
)

type UserRepository struct {
	*transactor
}

func NewUserRepository(dbGrm *gorm.DB) *UserRepository {
	return &UserRepository{
		transactor: NewTransactor(dbGrm),
	}
}

func (u *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	return u.getDB(ctx).Create(user).Error
}

func (u *UserRepository) GetAllUsers(ctx context.Context, req dto.PaginationReq) ([]*model.User, error) {
	var users []*model.User
	err := u.db.WithContext(ctx).
		Select("user_id", "user_email", "user_full_name", "user_avatar", "user_status", "user_created_at").
		Limit(req.GetLimit()).
		Offset(req.GetOffset()).
		Find(&users).Error

	if err != nil {
		return nil, err
	}
	return users, nil
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
	return u.getDB(ctx).Create(session).Error
}

func (u *UserRepository) RevokeSession(ctx context.Context, id uuid.UUID) error {
	return u.db.WithContext(ctx).Model(&model.Session{}).Where("sess_id = ?", id).Update("sess_is_blocked", true).Error
}
