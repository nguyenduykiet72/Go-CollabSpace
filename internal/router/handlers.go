package router

import "Go-CollabSpace/internal/controller"

type AppHandlers struct {
	AuthController      *controller.AuthController
	UserController      *controller.UserController
	WorkspaceController *controller.WorkspaceController
	DocumentController  *controller.DocumentController
	StorageController   *controller.StorageController
}
