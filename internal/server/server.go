package server

import (
	"Go-CollabSpace/config"
	"Go-CollabSpace/internal/common/infrastructure"
	"Go-CollabSpace/internal/common/token"
	"Go-CollabSpace/internal/controller"
	"Go-CollabSpace/internal/middleware"
	"Go-CollabSpace/internal/realtime"
	"Go-CollabSpace/internal/repository"
	"Go-CollabSpace/internal/router"
	"Go-CollabSpace/internal/service"
	"Go-CollabSpace/internal/worker"
	"Go-CollabSpace/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	cfg    *config.Config
	db     *gorm.DB
	engine *gin.Engine

	hub             *realtime.Hub
	redisClient     *infrastructure.RedisClient
	taskProcessor   worker.TaskProcessor
	taskDistributor worker.TaskDistributor
}

func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	return &Server{
		cfg: cfg,
		db:  db,
	}
}

func (s *Server) InitEngine() {
	gin.SetMode(resolveGinMode(s.cfg.Server.Mode))
	r := gin.New()

	tokenProvider := token.NewJWTProvider(token.ConfigToken(s.cfg.JWT))
	transactor := repository.NewTransactor(s.db)

	// -- Repositories --
	authRepo := repository.NewAuthRepository(s.db)
	userRepo := repository.NewUserRepository(s.db)
	workspaceRepo := repository.NewWorkspaceRepository(s.db)
	documentRepo := repository.NewDocumentRepository(s.db)

	redisClient, err := infrastructure.NewRedisClient(s.cfg.Redis)
	if err != nil {
		logger.Log.Info("Failed to connect to Redis", zap.Error(err))
	}
	s.redisClient = redisClient

	redisOpts := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", s.cfg.Redis.Host, s.cfg.Redis.Port),
		Password: s.cfg.Redis.Password,
		DB:       0,
	}

	emailSender := infrastructure.NewResendEmailSender(s.cfg.ResendEmail.ResendAPIKey, s.cfg.ResendEmail.FromEmail)

	taskDistributor := worker.NewRedisTaskDistributor(redisOpts)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpts, s.db, emailSender)
	s.taskDistributor = taskDistributor
	s.taskProcessor = taskProcessor
	go func() {
		logger.Log.Info("Starting task processor")
		if err := taskProcessor.Start(); err != nil {
			logger.Log.Fatal("Failed to start task processor", zap.Error(err))
		}
	}()

	s.hub = realtime.NewHub(documentRepo, redisClient)
	go s.hub.Run()
	wsController := controller.NewWsController(s.hub, tokenProvider, s.db, s.cfg.Server.AllowedOrigins)
	r.GET("/ws/*any", wsController.HandleWS)

	corsConfig := cors.Config{
		AllowOrigins:     s.cfg.Server.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))
	r.Use(gin.Recovery(), middleware.ErrorHandler())

	ctx := context.Background()
	s3Client, err := infrastructure.NewS3Client(ctx, s.cfg.AWS)
	if err != nil {
		logger.Log.Fatal("Failed to connect to AWS S3 ", zap.Error(err))
	}

	storageService := service.NewStorageService(s.db, s3Client, s.cfg.AWS.BucketName)
	storageController := controller.NewStorageController(storageService)

	googleProvider := infrastructure.NewGoogleOAuth(
		s.cfg.OAuthGoogleConfig.ClientID,
		s.cfg.OAuthGoogleConfig.ClientSecret,
		s.cfg.OAuthGoogleConfig.RedirectURL,
	)

	oauthProviders := map[string]infrastructure.OAuthProvider{
		"google": googleProvider,
		// Future providers can be added here like "github": githubProvider, etc.
	}

	// -- Auth Module --
	authService := service.NewAuthService(authRepo, userRepo, tokenProvider, transactor, taskDistributor, oauthProviders)
	authController := controller.NewAuthController(authService)

	// -- User Module --
	userService := service.NewUserService(userRepo, transactor)
	userController := controller.NewUserController(userService)

	// -- Workspace Module --
	workspaceService := service.NewWorkspaceService(workspaceRepo)
	workspaceController := controller.NewWorkspaceController(workspaceService)

	// -- Document Module --
	documentService := service.NewDocumentService(documentRepo, workspaceRepo, taskDistributor)
	documentController := controller.NewDocumentController(documentService)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	rateLimiter := middleware.NewRateLimiter(redisClient.Client)

	handlers := router.AppHandlers{
		AuthController:      authController,
		UserController:      userController,
		WorkspaceController: workspaceController,
		DocumentController:  documentController,
		StorageController:   storageController,
		RateLimiter:         rateLimiter,
	}

	router.SetUpRoutes(r, handlers, tokenProvider, s.db)

	s.engine = r
}

func (s *Server) Run(ctx context.Context) error {
	s.InitEngine()

	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	logger.Log.Info("Server is running", zap.String("addr", addr))

	srv := &http.Server{
		Addr:              addr,
		Handler:           s.engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		logger.Log.Info("Shutdown signal received, draining...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.shutdown(shutdownCtx, srv)
}

func (s *Server) shutdown(ctx context.Context, srv *http.Server) error {
	var firstErr error

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("HTTP server shutdown error", zap.Error(err))
		firstErr = err
	}

	if s.hub != nil {
		s.hub.Stop()
	}

	if s.taskProcessor != nil {
		s.taskProcessor.Shutdown()
	}

	if s.taskDistributor != nil {
		if err := s.taskDistributor.Close(); err != nil {
			logger.Log.Error("Task distributor close error", zap.Error(err))
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	if s.redisClient != nil {
		if err := s.redisClient.Close(); err != nil {
			logger.Log.Error("Redis close error", zap.Error(err))
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	if sqlDB, err := s.db.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			logger.Log.Error("DB close error", zap.Error(err))
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	logger.Log.Info("Graceful shutdown complete")
	return firstErr
}
