package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("requestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func Logger(logger interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		requestID := c.GetString("requestID")

		// Log using structured logger
		// logger.Info("HTTP Request",
		// 	"method", method,
		// 	"path", path,
		// 	"status", statusCode,
		// 	"latency", latency,
		// 	"ip", clientIP,
		// 	"request_id", requestID,
		// )

		_ = logger
		_ = method
		_ = path
		_ = statusCode
		_ = latency
		_ = clientIP
		_ = requestID
	}
}

func Recovery(logger interface{}) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		// logger.Error("Panic recovered", "error", err)
		c.JSON(500, gin.H{
			"success": false,
			"message": "Internal server error",
			"error":   "An unexpected error occurred",
		})
		c.Abort()

		_ = logger
	})
}
