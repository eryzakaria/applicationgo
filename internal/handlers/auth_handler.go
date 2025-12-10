package handlers

import (
	"net/http"

	"suitemedia/internal/models"
	"suitemedia/internal/service"
	"suitemedia/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "Registration data"
// @Success 201 {object} response.Response{data=models.AuthResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	authResp, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrUserEmailExists {
			response.Error(c, http.StatusConflict, "Email already exists", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to register user", err)
		return
	}

	response.Success(c, authResp, http.StatusCreated)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return access token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=models.AuthResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	authResp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			response.Error(c, http.StatusUnauthorized, "Invalid email or password", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to login", err)
		return
	}

	response.Success(c, authResp)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.Response{data=models.AuthResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	authResp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if err == service.ErrInvalidToken {
			response.Error(c, http.StatusUnauthorized, "Invalid refresh token", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to refresh token", err)
		return
	}

	response.Success(c, authResp)
}
