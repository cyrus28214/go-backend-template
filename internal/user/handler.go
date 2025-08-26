package user

import (
	"backend/pkg/apperror"
	"backend/pkg/contextkey"
	"backend/pkg/response"
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService IUserService
}

func NewUserHandler(userService IUserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	l := ctx.Value(contextkey.LoggerKey).(*slog.Logger)
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		l.Info("invalid request", "error", err)
		response.Error(c, apperror.ErrBadRequest.Wrap(err))
		return
	}

	token, err := h.userService.Login(ctx, &req)
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		response.Error(c, appErr)
		return
	} else {
		l.Error("unexpected error", "error", err)
	}

	response.Success(c, gin.H{"token": token})
}
