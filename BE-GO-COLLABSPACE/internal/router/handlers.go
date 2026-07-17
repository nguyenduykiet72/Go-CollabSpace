package router

import (
	"Go-CollabSpace/internal/controller"
	"Go-CollabSpace/internal/middleware"
)

type AppHandlers struct {
	AuthController      *controller.AuthController
	UserController      *controller.UserController
	WorkspaceController *controller.WorkspaceController
	DocumentController  *controller.DocumentController
	StorageController   *controller.StorageController

	RateLimiter *middleware.RateLimiter
}
