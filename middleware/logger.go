package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LogMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		status := c.Writer.Status()
		level := zap.InfoLevel
		if status != 200 {
			level = zap.WarnLevel
		}
		logger.Log(level, "HTTP request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", status),
			zap.Int("port", 8080),
			zap.Int64("latency_ms", duration.Milliseconds()),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
