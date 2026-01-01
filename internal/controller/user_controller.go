package controller

import (
	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/internal/service"
	"Go-CollabSpace/pkg/httpx"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.IUserService
}

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(apperror.ErrBadRequest)
		return
	}

	resp, err := c.userService.Register(ctx.Request.Context(), req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated,
		httpx.Success(resp, "User created successfully"),
	)
}

func (c *UserController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(apperror.ErrBadRequest)
		return
	}

	userAgent := ctx.GetHeader("User-Agent")
	if userAgent == "" {
		userAgent = "unknown"
	}

	tokenResp, err := c.userService.Login(ctx.Request.Context(), req, userAgent)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, httpx.Success(tokenResp, "Login successful"))
}
