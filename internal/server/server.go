package server

import (
	"Go-CollabSpace/config"
	"Go-CollabSpace/internal/common/token"
	"Go-CollabSpace/internal/controller"
	"Go-CollabSpace/internal/middleware"
	"Go-CollabSpace/internal/repository"
	"Go-CollabSpace/internal/router"
	"Go-CollabSpace/internal/service"
	"Go-CollabSpace/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	cfg    *config.Config
	db     *gorm.DB
	engine *gin.Engine
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
	r.Use(gin.Recovery(), middleware.ErrorHandler())

	tokenProvider := token.NewJWTProvider(token.ConfigToken(s.cfg.JWT))

	userRepo := repository.NewUserRepository(s.db)
	userService := service.NewUserService(userRepo, tokenProvider)
	userController := controller.NewUserController(userService)

	handlers := router.AppHandlers{
		UserController: userController,
	}

	router.SetUpRoutes(r, handlers)

	s.engine = r
}

func (s *Server) Run() error {
	s.InitEngine()

	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	logger.Log.Info("Server is running", zap.String("addr", addr))

	srv := &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}

	return srv.ListenAndServe()
}
