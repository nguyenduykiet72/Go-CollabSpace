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
	hub    *realtime.Hub
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
		logger.Log.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	redisOpts := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", s.cfg.Redis.Host, s.cfg.Redis.Port),
		Password: s.cfg.Redis.Password,
		DB:       0,
	}

	//emailSender := infrastructure.NewSMTPEmailSender(
	//	s.cfg.SMTP.Host,
	//	s.cfg.SMTP.Port,
	//	s.cfg.SMTP.Username,
	//	s.cfg.SMTP.Password,
	//	s.cfg.SMTP.From,
	//)

	emailSender := infrastructure.NewResendEmailSender(s.cfg.ResendEmail.ResendAPIKey, s.cfg.ResendEmail.FromEmail)

	taskDistributor := worker.NewRedisTaskDistributor(redisOpts)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpts, s.db, emailSender)
	go func() {
		logger.Log.Info("Starting task distributor")
		if err := taskProcessor.Start(); err != nil {
			logger.Log.Fatal("Failed to start task processor", zap.Error(err))
		}
	}()

	s.hub = realtime.NewHub(documentRepo, redisClient)
	go s.hub.Run()
	wsController := controller.NewWsController(s.hub, tokenProvider, s.db)
	r.GET("/ws/*any", wsController.HandleWS)

	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
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

	// -- Auth Module --
	authService := service.NewAuthService(authRepo, userRepo, tokenProvider, transactor)
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

	handlers := router.AppHandlers{
		AuthController:      authController,
		UserController:      userController,
		WorkspaceController: workspaceController,
		DocumentController:  documentController,
		StorageController:   storageController,
	}

	router.SetUpRoutes(r, handlers, tokenProvider, s.db)

	s.engine = r
}

func (s *Server) Run() error {
	s.InitEngine()

	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	logger.Log.Info("Server is running", zap.String("addr", addr))

	srv := &http.Server{
		Addr:              addr,
		Handler:           s.engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return srv.ListenAndServe()
}
