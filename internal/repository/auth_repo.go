package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/model"
)

type AuthRepository struct {
	*transactor
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		transactor: NewTransactor(db),
	}
}

func (a *AuthRepository) CreatePasswordReset(ctx context.Context, reset *model.PasswordReset) error {
	return a.getDB(ctx).Create(reset).Error
}

func (a *AuthRepository) FindValidPasswordReset(ctx context.Context, tokenHash string) (*model.PasswordReset, error) {
	var resetData model.PasswordReset

	err := a.getDB(ctx).
		Where("pass_token_hash = ?", tokenHash).
		Where("pass_is_used = ?", false).
		Where("pass_expire_at > ?", time.Now()).
		First(&resetData).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}

	return &resetData, nil
}

func (a *AuthRepository) MarkTokenAsUsed(ctx context.Context, resetID uuid.UUID) error {
	return a.getDB(ctx).
		Model(&model.PasswordReset{}).
		Where("pass_id = ?", resetID).
		Update("pass_is_used", true).Error
}
