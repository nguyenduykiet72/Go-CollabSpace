package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/common/infrastructure"
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/internal/model"
	"Go-CollabSpace/internal/worker"
	"Go-CollabSpace/pkg/hash"
	"Go-CollabSpace/pkg/logger"
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
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
	FindActiveSessionByRefreshHash(ctx context.Context, tokenHash string) (*model.Session, error)
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, newHashedPassword string) error
	UpdateSocialAuth(ctx context.Context, userID uuid.UUID, provider string, socialID string, avatar string) error
}

type TokenProvider interface {
	GenerateAccessToken(userID uuid.UUID) (string, error)
	GenerateRefreshToken() (string, error)
}

type AuthService struct {
	authRepo        AuthRepo
	userRepo        UserAuthRepo
	tokenProvider   TokenProvider
	transactor      Transactor
	taskDistributor worker.TaskDistributor
	oauthProviders  map[string]infrastructure.OAuthProvider
}

func NewAuthService(authRepo AuthRepo, userRepo UserAuthRepo, tokenProvider TokenProvider, transactor Transactor, taskDistributor worker.TaskDistributor, providers map[string]infrastructure.OAuthProvider) *AuthService {
	return &AuthService{
		authRepo:        authRepo,
		userRepo:        userRepo,
		tokenProvider:   tokenProvider,
		transactor:      transactor,
		taskDistributor: taskDistributor,
		oauthProviders:  providers,
	}
}

func (a *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
	existingUser, err := a.userRepo.GetUserByEmail(ctx, req.Email)

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

	hashedPwdStr := string(hashPassword)

	newUser := &model.User{
		UserEmail:    req.Email,
		UserPassword: &hashedPwdStr,
		UserFullName: req.FullName,
		UserStatus:   "active",
	}

	if err := a.userRepo.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       newUser.UserID,
		Email:    newUser.UserEmail,
		FullName: newUser.UserFullName,
	}, nil
}

func (a *AuthService) Login(ctx context.Context, req dto.LoginRequest, userAgent string) (*dto.TokenResponse, error) {
	invalidCredErr := apperror.NewAppError(http.StatusBadRequest, "invalid email or password", "")

	user, err := a.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, invalidCredErr
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.UserPassword), []byte(req.Password)); err != nil {
		return nil, invalidCredErr
	}

	return a.issueTokensTx(ctx, user.UserID, userAgent)
}

// generateTokenPair generates a fresh access + refresh token pair without
// touching the database.
func (a *AuthService) generateTokenPair(userID uuid.UUID) (accessToken, refreshToken string, err error) {
	accessToken, err = a.tokenProvider.GenerateAccessToken(userID)
	if err != nil {
		return "", "", apperror.ErrInternal.WithRootErr(err)
	}
	refreshToken, err = a.tokenProvider.GenerateRefreshToken()
	if err != nil {
		return "", "", apperror.ErrInternal.WithRootErr(err)
	}
	return accessToken, refreshToken, nil
}

// persistSession stores a hashed refresh token for the given user. The caller
// owns the transaction (use ctx scoped by Transactor.WithinTransaction).
func (a *AuthService) persistSession(ctx context.Context, userID uuid.UUID, refreshToken, userAgent string) error {
	const refreshTokenTTL = 7 * 24 * time.Hour
	session := &model.Session{
		SessUserID:       userID,
		SessRefreshToken: hash.HashToken(refreshToken),
		SessUserAgent:    userAgent,
		SessExpireAt:     time.Now().Add(refreshTokenTTL),
	}
	if err := a.userRepo.CreateSession(ctx, session); err != nil {
		return apperror.ErrInternal.WithRootErr(err)
	}
	return nil
}

// issueTokensTx generates a token pair and persists the session inside a new
// transaction.
func (a *AuthService) issueTokensTx(ctx context.Context, userID uuid.UUID, userAgent string) (*dto.TokenResponse, error) {
	accessToken, refreshToken, err := a.generateTokenPair(userID)
	if err != nil {
		return nil, err
	}
	err = a.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		return a.persistSession(txCtx, userID, refreshToken, userAgent)
	})
	if err != nil {
		return nil, err
	}
	return &dto.TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// Refresh validates the supplied refresh token, rotates it (revokes the old
// session, issues a new one) and returns a fresh token pair. Re-use of an old
// or unknown refresh token returns ErrInvalidToken.
func (a *AuthService) Refresh(ctx context.Context, refreshToken, userAgent string) (*dto.TokenResponse, error) {
	if refreshToken == "" {
		return nil, apperror.ErrInvalidToken
	}

	tokenHash := hash.HashToken(refreshToken)
	session, err := a.userRepo.FindActiveSessionByRefreshHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, apperror.ErrInvalidToken
		}
		return nil, err
	}

	newAccess, newRefresh, err := a.generateTokenPair(session.SessUserID)
	if err != nil {
		return nil, err
	}

	err = a.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := a.userRepo.RevokeSession(txCtx, session.SessID); err != nil {
			return apperror.ErrInternal.WithRootErr(err)
		}
		return a.persistSession(txCtx, session.SessUserID, newRefresh, userAgent)
	})
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{AccessToken: newAccess, RefreshToken: newRefresh}, nil
}

func (a *AuthService) ForgotPassword(ctx context.Context, dto dto.ForgotPasswordRequest) error {
	user, err := a.userRepo.GetUserByEmail(ctx, dto.Email)
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

	if err := a.authRepo.CreatePasswordReset(ctx, resetData); err != nil {
		return err
	}

	payload := &worker.PayloadSendResetEmail{
		ToEmail:    dto.Email,
		ResetToken: rawToken,
	}

	// We intentionally don't surface enqueue errors to the caller (timing-attack
	// safe: callers should not learn whether the email exists). However we log
	// them so operators can spot a broken queue.
	if err := a.taskDistributor.DistributeTaskSendResetEmail(ctx, payload); err != nil {
		logger.Log.Error("failed to enqueue reset email task",
			zap.String("email", dto.Email),
			zap.Error(err),
		)
	}

	return nil
}

func (a *AuthService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	hashedToken := hash.HashToken(req.RawToken)

	resetData, err := a.authRepo.FindValidPasswordReset(ctx, hashedToken)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return apperror.ErrInvalidToken
		}
		return err
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperror.ErrInternal.WithRootErr(err)
	}

	return a.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := a.userRepo.UpdateUserPassword(txCtx, resetData.PassUserID, string(hashedPwd)); err != nil {
			return err
		}

		if err := a.authRepo.MarkTokenAsUsed(txCtx, resetData.PassID); err != nil {
			return err
		}

		return nil
	})
}

func (a *AuthService) LoginWithSocial(ctx context.Context, providerName string, code string, userAgent string) (*dto.TokenResponse, error) {
	provider, exists := a.oauthProviders[providerName]
	if !exists {
		return nil, apperror.ErrNotFound
	}

	profile, err := provider.GetProfile(ctx, code)
	if err != nil {
		return nil, apperror.ErrFetchProfileFailed
	}

	user, err := a.resolveSocialUser(ctx, profile, providerName)
	if err != nil {
		return nil, err
	}

	return a.issueTokensTx(ctx, user.UserID, userAgent)
}

func (a *AuthService) resolveSocialUser(ctx context.Context, profile *infrastructure.SocialProfile, providerName string) (*model.User, error) {
	existingUser, err := a.userRepo.GetUserByEmail(ctx, profile.Email)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			newUser := &model.User{
				UserEmail:    profile.Email,
				UserFullName: profile.Name,
				UserAvatar:   profile.Avatar,
				UserStatus:   "active",
				AuthProvider: providerName,
				SocialID:     &profile.ID,
			}
			if err := a.userRepo.CreateUser(ctx, newUser); err != nil {
				return nil, err
			}
			return newUser, nil
		}
		return nil, err
	}

	needUpdate := existingUser.AuthProvider != providerName ||
		existingUser.SocialID == nil ||
		*existingUser.SocialID != profile.ID

	if needUpdate {
		if err := a.userRepo.UpdateSocialAuth(ctx, existingUser.UserID, providerName, profile.ID, profile.Avatar); err != nil {
			// Login can still proceed with the stale profile, but log so we can
			// reconcile later.
			logger.Log.Warn("failed to refresh social auth on login",
				zap.String("user_id", existingUser.UserID.String()),
				zap.String("provider", providerName),
				zap.Error(err),
			)
		} else {
			existingUser.AuthProvider = providerName
			existingUser.SocialID = &profile.ID
			existingUser.UserAvatar = profile.Avatar
		}
	}

	return existingUser, nil
}
