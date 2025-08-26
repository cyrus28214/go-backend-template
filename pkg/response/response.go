package response

import (
	"backend/pkg/apperror"
	"backend/pkg/contextkey"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Response struct {
	Code    string     `json:"code"`
	Data    any        `json:"data,omitempty"`
	TraceID *uuid.UUID `json:"trace_id,omitempty"`
}

func Success(c *gin.Context, data any) {
	ctx := c.Request.Context()
	var traceID *uuid.UUID
	if traceIDValue := ctx.Value(contextkey.TraceIDKey); traceIDValue != nil {
		if id, ok := traceIDValue.(uuid.UUID); ok {
			traceID = &id
		}
	}

	c.JSON(http.StatusOK, Response{
		Code:    "OK",
		Data:    data,
		TraceID: traceID,
	})
}

func Error(c *gin.Context, err *apperror.AppError) {
	// 安全检查：如果 err 为 nil，使用默认的内部服务器错误
	if err == nil {
		err = apperror.ErrInternal
	}

	ctx := c.Request.Context()
	var traceID *uuid.UUID
	if traceIDValue := ctx.Value(contextkey.TraceIDKey); traceIDValue != nil {
		if id, ok := traceIDValue.(uuid.UUID); ok {
			traceID = &id
		}
	}

	c.JSON(err.HTTPStatus, Response{
		Code:    err.Code,
		Data:    err.Data,
		TraceID: traceID,
	})
}
