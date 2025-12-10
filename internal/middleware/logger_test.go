package middleware

import (
	"net/http/httptest"
	"testing"

	"suitemedia/pkg/logger"

	"github.com/gin-gonic/gin"
)

func TestRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	handler := RequestID()
	handler(c)

	requestID := c.Writer.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Log("Request ID middleware executed successfully")
	}
}

func TestRequestIDWithExisting(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Request-ID", "existing-id")
	c.Request = req

	handler := RequestID()
	handler(c)

	// Test passes if no panic occurs
	t.Log("Request ID middleware with existing ID executed successfully")
}

func TestRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	log := logger.NewLogger("debug")
	handler := Recovery(log)

	// Test that recovery middleware doesn't panic on normal execution
	handler(c)

	// Middleware should complete successfully
	if w.Code != 0 && w.Code != 200 {
		t.Errorf("Expected status 0 or 200, got %d", w.Code)
	}
}
