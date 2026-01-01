package service

import (
	"Go-CollabSpace/internal/common/token"
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/internal/model"
	"Go-CollabSpace/internal/repository"
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUserService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req dto.LoginRequest, userAgent string) (*dto.TokenResponse, error)
}

type UserService struct {
	userRepo      repository.IUserRepository
	tokenProvider token.ITokenProvider
}

func NewUserService(userRepo repository.IUserRepository, tokenProvider token.ITokenProvider) IUserService {
	return &UserService{userRepo: userRepo, tokenProvider: tokenProvider}
}

func (u *UserService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
	existingUser, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
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
	user, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := u.tokenProvider.GenerateAccessToken(user.UserID, "Viewer")
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.tokenProvider.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshTokenTTL := 7 * 24 * time.Hour

	session := &model.Session{
		SessUserID:       user.UserID,
		SessRefreshToken: refreshToken,
		SessUserAgent:    userAgent,
		SessExpireAt:     time.Now().Add(refreshTokenTTL),
	}

	if err := u.userRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
