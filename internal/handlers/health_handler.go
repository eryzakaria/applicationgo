package handlers

import (
	"database/sql"
	"net/http"

	"suitemedia/pkg/redis"
	"suitemedia/pkg/response"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	db    *sql.DB
	redis *redis.Client
}

func NewHealthHandler(db *sql.DB, redis *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:    db,
		redis: redis,
	}
}

// Health godoc
// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Produce json
// @Success 200 {object} response.Response
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	response.Success(c, gin.H{
		"status":  "healthy",
		"service": "suitemedia-api",
	})
}

// Ready godoc
// @Summary Readiness check
// @Description Check if the service is ready to accept requests
// @Tags health
// @Produce json
// @Success 200 {object} response.Response
// @Failure 503 {object} response.Response
// @Router /ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	checks := map[string]bool{
		"database": h.checkDatabase(),
		"redis":    h.checkRedis(),
	}

	allHealthy := true
	for _, healthy := range checks {
		if !healthy {
			allHealthy = false
			break
		}
	}

	if !allHealthy {
		response.Error(c, http.StatusServiceUnavailable, "Service not ready", nil)
		return
	}

	response.Success(c, gin.H{
		"status": "ready",
		"checks": checks,
	})
}

func (h *HealthHandler) checkDatabase() bool {
	if h.db == nil {
		return false
	}
	return h.db.Ping() == nil
}

func (h *HealthHandler) checkRedis() bool {
	if h.redis == nil {
		return false
	}
	return h.redis.Ping() == nil
}
