package apperror

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code       string
	HTTPStatus int
	Data       any
	cause      error
}

func (e *AppError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("code: %s, data: %v, cause: %v", e.Code, e.Data, e.cause)
	}
	return fmt.Sprintf("code: %s, data: %v", e.Code, e.Data)
}

func (e *AppError) Unwrap() error {
	return e.cause
}

func (e *AppError) Wrap(cause error) *AppError {
	newErr := *e
	newErr.cause = cause
	return &newErr
}

func (e *AppError) WithData(data any) *AppError {
	newErr := *e
	newErr.Data = data
	return &newErr
}

func NewAppError(code string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		HTTPStatus: httpStatus,
	}
}

var (
	ErrInternal   = NewAppError("INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
	ErrBadRequest = NewAppError("BAD_REQUEST", http.StatusBadRequest)
)
