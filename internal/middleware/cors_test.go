package middleware

import (
	"net/http/httptest"
	"testing"

	"suitemedia/config"

	"github.com/gin-gonic/gin"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.CORSConfig{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("OPTIONS", "/", nil)
	c.Request.Header.Set("Origin", "http://localhost:3000")

	handler := CORS(cfg)
	handler(c)

	// Test passes if no panic occurs
	t.Log("CORS middleware executed successfully")
}
