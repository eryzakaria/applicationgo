package handlers

import (
	"net/http"

	"suitemedia/internal/models"
	"suitemedia/internal/service"
	"suitemedia/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// List godoc
// @Summary List users
// @Description Get list of users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=response.PaginatedData}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users [get]
func (h *UserHandler) List(c *gin.Context) {
	var params models.ListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid query parameters", err)
		return
	}

	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	users, total, err := h.userService.List(c.Request.Context(), params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch users", err)
		return
	}

	response.SuccessPaginated(c, users, params.Page, params.Limit, total)
}

// GetByID godoc
// @Summary Get user by ID
// @Description Get user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=models.UserResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrUserNotFound {
			response.Error(c, http.StatusNotFound, "User not found", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to fetch user", err)
		return
	}

	response.Success(c, user)
}

// Create godoc
// @Summary Create new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User data"
// @Security BearerAuth
// @Success 201 {object} response.Response{data=models.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	user, err := h.userService.Create(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrUserEmailExists {
			response.Error(c, http.StatusConflict, "Email already exists", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	response.Success(c, user, http.StatusCreated)
}

// Update godoc
// @Summary Update user
// @Description Update user details
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UpdateUserRequest true "User data"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=models.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	user, err := h.userService.Update(c.Request.Context(), id, req)
	if err != nil {
		if err == service.ErrUserNotFound {
			response.Error(c, http.StatusNotFound, "User not found", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	response.Success(c, user)
}

// Delete godoc
// @Summary Delete user
// @Description Delete user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.userService.Delete(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrUserNotFound {
			response.Error(c, http.StatusNotFound, "User not found", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}

	response.Success(c, gin.H{"message": "User deleted successfully"})
}

// GetProfile godoc
// @Summary Get current user profile
// @Description Get authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=models.UserResponse}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("userID")

	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch profile", err)
		return
	}

	response.Success(c, user)
}

// UpdateProfile godoc
// @Summary Update current user profile
// @Description Update authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.UpdateUserRequest true "User data"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=models.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/me [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("userID")

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	user, err := h.userService.Update(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update profile", err)
		return
	}

	response.Success(c, user)
}
