package service

import (
	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/internal/model"
	"Go-CollabSpace/internal/worker"
	"Go-CollabSpace/pkg/hash"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo interface {
	CreatePasswordReset(ctx context.Context, reset *model.PasswordReset) error
	FindValidPasswordReset(ctx context.Context, tokenHash string) (*model.PasswordReset, error)
	MarkTokenAsUsed(ctx context.Context, resetID uuid.UUID) error
}

type UserAuthRepo interface {
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	CreateSession(ctx context.Context, session *model.Session) error
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, newHashedPassword string) error
}

type TokenProvider interface {
	GenerateAccessToken(userId uuid.UUID, role string) (string, error)
	GenerateRefreshToken() (string, error)
}

type AuthService struct {
	authRepo        AuthRepo
	userRepo        UserAuthRepo
	tokenProvider   TokenProvider
	transactor      Transactor
	taskDistributor worker.TaskDistributor
}

func NewAuthService(authRepo AuthRepo, userRepo UserAuthRepo, tokenProvider TokenProvider, transactor Transactor) *AuthService {
	return &AuthService{authRepo: authRepo, userRepo: userRepo, tokenProvider: tokenProvider, transactor: transactor}
}

func (u *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
	existingUser, err := u.userRepo.GetUserByEmail(ctx, req.Email)

	if err != nil && !errors.Is(err, apperror.ErrNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, apperror.ErrEmailExists
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		UserEmail:    req.Email,
		UserPassword: string(hashPassword),
		UserFullName: req.FullName,
		UserStatus:   "active",
	}

	if err := u.userRepo.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       newUser.UserID,
		Email:    newUser.UserEmail,
		FullName: newUser.UserFullName,
	}, nil
}

func (u *AuthService) Login(ctx context.Context, req dto.LoginRequest, userAgent string) (*dto.TokenResponse, error) {
	invalidCredErr := apperror.NewAppError(http.StatusBadRequest, "invalid email or password", "")

	user, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, invalidCredErr
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(req.Password)); err != nil {
		return nil, invalidCredErr
	}

	accessToken, err := u.tokenProvider.GenerateAccessToken(user.UserID, "Viewer")
	if err != nil {
		return nil, apperror.ErrInternal.WithRootErr(err)
	}

	refreshToken, err := u.tokenProvider.GenerateRefreshToken()
	if err != nil {
		return nil, apperror.ErrInternal.WithRootErr(err)
	}

	refreshTokenTTL := 7 * 24 * time.Hour

	err = u.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		session := &model.Session{
			SessUserID:       user.UserID,
			SessRefreshToken: refreshToken,
			SessUserAgent:    userAgent,
			SessExpireAt:     time.Now().Add(refreshTokenTTL),
		}

		if err := u.userRepo.CreateSession(txCtx, session); err != nil {
			return apperror.ErrInternal.WithRootErr(err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *AuthService) ForgotPassword(ctx context.Context, dto dto.ForgotPasswordRequest) error {
	user, err := u.userRepo.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil // Don't reveal if email exists
		}
		return err
	}

	rawToken := uuid.New().String() + uuid.New().String()
	hashedToken := hash.HashToken(rawToken)

	resetData := &model.PasswordReset{
		PassUserID:    user.UserID,
		PassTokenHash: hashedToken,
		PassExpireAt:  time.Now().Add(15 * time.Minute),
	}

	if err := u.authRepo.CreatePasswordReset(ctx, resetData); err != nil {
		return err
	}

	payload := &worker.PayloadSendResetEmail{
		ToEmail:    dto.Email,
		ResetToken: rawToken,
	}

	_ = u.taskDistributor.DistributeTaskSendResetEmail(ctx, payload)

	return nil
}

func (u *AuthService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	hashedToken := hash.HashToken(req.RawToken)

	resetData, err := u.authRepo.FindValidPasswordReset(ctx, hashedToken)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return apperror.ErrInvalidToken
		}
		return err
	}

	return u.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)

		if err := u.userRepo.UpdateUserPassword(txCtx, resetData.PassUserID, string(hashedPwd)); err != nil {
			return err
		}

		if err := u.authRepo.MarkTokenAsUsed(txCtx, resetData.PassID); err != nil {
			return err
		}

		return nil
	})
}
