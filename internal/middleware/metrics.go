package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// Calculate request duration
		duration := time.Since(start)

		// Log metrics (in production, send to Prometheus)
		// Example: prometheus.ObserveRequestDuration(c.Request.Method, c.Request.URL.Path, duration)
		_ = duration // Avoid unused variable warning
	}
}
