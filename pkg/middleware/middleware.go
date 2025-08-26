package middleware

import (
	"backend/pkg/apperror"
	"backend/pkg/contextkey"
	"backend/pkg/response"
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoggerMiddleware(l *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		traceID, _ := uuid.NewV7()

		logger := l.With(
			"trace_id", traceID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"params", c.ClientIP(),
			// "user_agent", c.Request.UserAgent(),
		)

		logger.Info("request received")

		ctx := context.WithValue(c.Request.Context(), contextkey.TraceIDKey, traceID)
		ctx = context.WithValue(ctx, contextkey.LoggerKey, logger)

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		latency := time.Since(start)

		logger.Info("request completed", "status", c.Writer.Status(), "latency", latency)
	}
}

func RecoveryMiddleware(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			l := c.Request.Context().Value(contextkey.LoggerKey).(*slog.Logger)
			l.Error("panic recovered", "error", r, "stack", string(debug.Stack()))
			response.Error(c, apperror.ErrInternal)
		}
	}()
	c.Next()
}
