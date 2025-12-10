package middleware

import (
	"github.com/gin-gonic/gin"
)

var (
	requestCounter  = 0
	requestDuration = 0.0
)

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Increment request counter
		requestCounter++

		c.Next()

		// Record metrics (placeholder)
		// In production, use Prometheus metrics
	}
}
