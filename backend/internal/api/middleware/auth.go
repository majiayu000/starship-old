package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/majiayu000/cc-starship/internal/core/domain"
	"github.com/majiayu000/cc-starship/internal/core/ports"
	"github.com/majiayu000/cc-starship/pkg/errors"
	"github.com/majiayu000/cc-starship/pkg/utils"
)

// Auth creates a middleware that authenticates requests
func Auth(authService ports.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, errors.NewUnauthorized("Authorization header required", nil))
			c.Abort()
			return
		}

		// Check if the Authorization header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, errors.NewUnauthorized("Invalid authorization format", nil))
			c.Abort()
			return
		}

		// Extract the token
		token := parts[1]

		// Validate token
		user, err := authService.ValidateToken(c, token)
		if err != nil {
			utils.ErrorResponse(c, errors.NewUnauthorized("Invalid or expired token", err))
			c.Abort()
			return
		}

		// Add user to context
		c.Set("user", user)

		// Call the next handler
		c.Next()
	}
}

// RequireRole creates a middleware that requires a specific role
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context
		user, exists := c.Get("user")
		if !exists {
			utils.ErrorResponse(c, errors.NewUnauthorized("User not authenticated", nil))
			c.Abort()
			return
		}

		// Check if user has required role
		userObj, ok := user.(*domain.User)
		if !ok {
			utils.ErrorResponse(c, errors.NewUnauthorized("Invalid user type", nil))
			c.Abort()
			return
		}

		if userObj.Role != requiredRole && userObj.Role != "admin" {
			utils.ErrorResponse(c, errors.NewForbidden("Insufficient permissions", nil))
			c.Abort()
			return
		}

		// Call the next handler
		c.Next()
	}
}
