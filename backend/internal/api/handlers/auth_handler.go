package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/majiayu000/cc-starship/internal/core/ports"
	"github.com/majiayu000/cc-starship/pkg/errors"
	"github.com/majiayu000/cc-starship/pkg/logger"
	"github.com/majiayu000/cc-starship/pkg/utils"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService ports.AuthService
	logger      *logger.Logger
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService ports.AuthService, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the response for authentication operations
type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// AuthStatusResponse reports whether the server requires authentication.
type AuthStatusResponse struct {
	Enabled bool `json:"enabled"`
}

// AuthStatus returns whether auth middleware is active for protected APIs.
// Always public so the frontend can load review data when auth is disabled.
func AuthStatus(enabled bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.JSONResponse(c, http.StatusOK, AuthStatusResponse{Enabled: enabled})
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	// Parse request body
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequest("Invalid request body", err))
		return
	}

	// Register user
	user, err := h.authService.Register(c, req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Generate token
	token, err := h.authService.Login(c, req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(c, errors.NewInternal("Failed to generate token", err))
		return
	}

	// Create response
	response := AuthResponse{
		Token: token,
		User: map[string]interface{}{
			"id":        user.ID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"role":      user.Role,
		},
	}

	// Send response
	utils.JSONResponse(c, http.StatusCreated, response)
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	// Parse request body
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequest("Invalid request body", err))
		return
	}

	// Authenticate user
	token, err := h.authService.Login(c, req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Get user
	user, err := h.authService.ValidateToken(c, token)
	if err != nil {
		utils.ErrorResponse(c, errors.NewInternal("Failed to validate token", err))
		return
	}

	// Create response
	response := AuthResponse{
		Token: token,
		User: map[string]interface{}{
			"id":        user.ID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"role":      user.Role,
		},
	}

	// Send response
	utils.JSONResponse(c, http.StatusOK, response)
}
