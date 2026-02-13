package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"Go-CollabSpace/config"
	"Go-CollabSpace/internal/common/token"
	"Go-CollabSpace/internal/controller"
	"Go-CollabSpace/internal/middleware"
	"Go-CollabSpace/internal/realtime"
	"Go-CollabSpace/internal/repository"
	"Go-CollabSpace/internal/router"
	"Go-CollabSpace/internal/service"
	"Go-CollabSpace/pkg/logger"
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

	s.hub = realtime.NewHub(documentRepo)
	go s.hub.Run()
	wsController := controller.NewWsController(s.hub, tokenProvider)
	r.GET("/ws/*any", wsController.HandleWS)

	r.Use(gin.Recovery(), middleware.ErrorHandler())

	// -- User Module --
	userService := service.NewUserService(userRepo, tokenProvider, transactor)
	userController := controller.NewUserController(userService)

	// -- Workspace Module --
	workspaceService := service.NewWorkspaceService(workspaceRepo)
	workspaceController := controller.NewWorkspaceController(workspaceService)

	// -- Document Module --
	documentService := service.NewDocumentService(documentRepo, workspaceRepo)
	documentController := controller.NewDocumentController(documentService)

	handlers := router.AppHandlers{
		UserController:      userController,
		WorkspaceController: workspaceController,
		DocumentController:  documentController,
	}

	router.SetUpRoutes(r, handlers, tokenProvider)

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
