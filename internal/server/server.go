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
	"Go-CollabSpace/pkg/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	userRepo := repository.NewUserRepository(s.db)
	workspaceRepo := repository.NewWorkspaceRepository(s.db)
	documentRepo := repository.NewDocumentRepository(s.db)

	redisClient, err := infrastructure.NewRedisClient(s.cfg.Redis)
	if err != nil {
		logger.Log.Info("Failed to connect to Redis", zap.Error(err))
	}

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

	// -- Auth Module --
	authService := service.NewAuthService(userRepo, tokenProvider, transactor)
	authController := controller.NewAuthController(authService)

	// -- User Module --
	userService := service.NewUserService(userRepo, transactor)
	userController := controller.NewUserController(userService)

	// -- Workspace Module --
	workspaceService := service.NewWorkspaceService(workspaceRepo)
	workspaceController := controller.NewWorkspaceController(workspaceService)

	// -- Document Module --
	documentService := service.NewDocumentService(documentRepo, workspaceRepo)
	documentController := controller.NewDocumentController(documentService)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	handlers := router.AppHandlers{
		AuthController:      authController,
		UserController:      userController,
		WorkspaceController: workspaceController,
		DocumentController:  documentController,
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
