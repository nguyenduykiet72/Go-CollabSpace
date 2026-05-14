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
	err := u.getDB(ctx).
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
	err := u.getDB(ctx).Where("user_email = ?", email).First(&user).Error
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
	err := u.getDB(ctx).First(&user, "user_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) CreateSession(ctx context.Context, session *model.Session) error {
	return u.getDB(ctx).Create(session).Error
}

func (u *UserRepository) RevokeSession(ctx context.Context, id uuid.UUID) error {
	return u.getDB(ctx).Model(&model.Session{}).Where("sess_id = ?", id).Update("sess_is_blocked", true).Error
}

// FindActiveSessionByRefreshHash returns the non-blocked, non-expired session
// matching the given hashed refresh token, if any.
func (u *UserRepository) FindActiveSessionByRefreshHash(ctx context.Context, tokenHash string) (*model.Session, error) {
	var session model.Session
	err := u.getDB(ctx).
		Where("sess_refresh_token = ?", tokenHash).
		Where("sess_is_blocked = ?", false).
		Where("sess_expire_at > NOW()").
		First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return &session, nil
}

func (u *UserRepository) UpdateUserPassword(ctx context.Context, userID uuid.UUID, newHashedPassword string) error {
	result := u.getDB(ctx).
		Model(&model.User{}).
		Where("user_id = ?", userID).
		Update("user_password", newHashedPassword)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperror.ErrNotFound
	}

	return nil
}

func (u *UserRepository) UpdateSocialAuth(ctx context.Context, userID uuid.UUID, provider string, socialID string, avatar string) error {
	updates := map[string]interface{}{
		"auth_provider": provider,
		"social_id":     socialID,
		"user_avatar":   avatar,
	}

	result := u.getDB(ctx).
		Model(&model.User{}).
		Where("user_id = ?", userID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperror.ErrNotFound
	}

	return nil
}
