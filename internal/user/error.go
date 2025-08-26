package user

import (
	"backend/pkg/apperror"
	"net/http"
)

var (
	ErrUserNotFound       = apperror.NewAppError("USER_NOT_FOUND", http.StatusNotFound)
	ErrInvalidCredentials = apperror.NewAppError("INVALID_CREDENTIALS", http.StatusUnauthorized)
)
