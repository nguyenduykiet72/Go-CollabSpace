package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/internal/model"
)

type UserRepo interface {
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	CreateSession(ctx context.Context, session *model.Session) error
	GetAllUsers(ctx context.Context, req dto.PaginationReq) ([]*model.User, error)
}

type TokenProvider interface {
	GenerateAccessToken(userId uuid.UUID, role string) (string, error)
	GenerateRefreshToken() (string, error)
}

type UserService struct {
	userRepo      UserRepo
	tokenProvider TokenProvider
	transactor    Transactor
}

func NewUserService(userRepo UserRepo, tokenProvider TokenProvider, transactor Transactor) *UserService {
	return &UserService{userRepo: userRepo, tokenProvider: tokenProvider, transactor: transactor}
}

func (u *UserService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
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

func (u *UserService) Login(ctx context.Context, req dto.LoginRequest, userAgent string) (*dto.TokenResponse, error) {
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
