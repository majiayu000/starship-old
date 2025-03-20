package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/majiayu000/cc-starship/pkg/errors"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// JSON sends a JSON response using Gin
func JSONResponse(c *gin.Context, status int, data interface{}) {
	// Create success response
	response := Response{
		Success: status >= 200 && status < 400,
		Data:    data,
	}

	// Send response
	c.JSON(status, response)
}

// ErrorResponse sends an error response using Gin
func ErrorResponse(c *gin.Context, err error) {
	var status int
	var errorResponse interface{}

	// Check if it's an application error
	if appErr, ok := errors.IsAppError(err); ok {
		status = appErr.StatusCode()
		errorResponse = map[string]string{
			"message": appErr.Message,
		}

		// Include details if available
		if appErr.Details != "" {
			errorResponse = map[string]string{
				"message": appErr.Message,
				"details": appErr.Details,
			}
		}
	} else {
		// For generic errors, use internal server error
		status = http.StatusInternalServerError
		errorResponse = map[string]string{
			"message": "An internal server error occurred",
		}
	}

	// Create error response
	response := Response{
		Success: false,
		Error:   errorResponse,
	}

	// Send response
	c.JSON(status, response)
}
