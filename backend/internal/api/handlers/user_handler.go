package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/majiayu000/cc-starship/internal/api/middleware"
	"github.com/majiayu000/cc-starship/internal/core/domain"
	"github.com/majiayu000/cc-starship/internal/core/ports"
	"github.com/majiayu000/cc-starship/pkg/errors"
	"github.com/majiayu000/cc-starship/pkg/logger"
	"github.com/majiayu000/cc-starship/pkg/utils"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userService ports.UserService
	logger      *logger.Logger
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService ports.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// UserRequest represents the request body for user operations
type UserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Role      string `json:"role"`
	Active    bool   `json:"active"`
}

// GetUsers handles the request to get all users
func (h *UserHandler) GetUsers(c *gin.Context) {
	// Get users from service
	users, err := h.userService.GetUsers(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Send response
	utils.JSONResponse(c, http.StatusOK, users)
}

// GetUser handles the request to get a user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	// Get user ID from request path parameters
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, errors.NewBadRequest("User ID is required", nil))
		return
	}

	// Get user from service
	user, err := h.userService.GetUserByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Send response
	utils.JSONResponse(c, http.StatusOK, user)
}

// CreateUser handles the request to create a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	// Parse request body
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequest("Invalid request body", err))
		return
	}

	// Create user
	user := domain.NewUser(req.Email, "", req.FirstName, req.LastName)
	if req.Role != "" {
		user.Role = req.Role
	}
	user.Active = req.Active

	// Save user
	if err := h.userService.CreateUser(c, user); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Send response
	utils.JSONResponse(c, http.StatusCreated, user)
}

// UpdateUser handles the request to update a user.
// Callers must be the same user (self profile edit) or an admin.
// Non-admin callers may update profile fields only; Role and Active are ignored.
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, errors.NewBadRequest("User ID is required", nil))
		return
	}

	caller, ok := middleware.CurrentUser(c)
	if !ok {
		utils.ErrorResponse(c, errors.NewUnauthorized("User not authenticated", nil))
		return
	}

	isAdmin := caller.Role == "admin"
	if !isAdmin && caller.ID != id {
		utils.ErrorResponse(c, errors.NewForbidden("Insufficient permissions", nil))
		return
	}

	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequest("Invalid request body", err))
		return
	}

	user, err := h.userService.GetUserByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	user.Email = req.Email
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	if isAdmin {
		if req.Role != "" && req.Role != user.Role {
			// Refuse demoting the last remaining administrator; otherwise no
			// caller can grant admin again through the API.
			if user.Role == "admin" && req.Role != "admin" {
				users, listErr := h.userService.GetUsers(c)
				if listErr != nil {
					utils.ErrorResponse(c, listErr)
					return
				}
				adminCount := 0
				for _, existing := range users {
					if existing != nil && existing.Role == "admin" {
						adminCount++
					}
				}
				if adminCount <= 1 {
					utils.ErrorResponse(c, errors.NewForbidden("Cannot demote the sole administrator", nil))
					return
				}
			}
			user.Role = req.Role
		}
		user.Active = req.Active
	}

	if err := h.userService.UpdateUser(c, user); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.JSONResponse(c, http.StatusOK, user)
}

// DeleteUser handles the request to delete a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Get user ID from request path parameters
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, errors.NewBadRequest("User ID is required", nil))
		return
	}

	// Delete user
	if err := h.userService.DeleteUser(c, id); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Send response
	utils.JSONResponse(c, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
