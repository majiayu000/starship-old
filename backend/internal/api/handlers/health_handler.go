package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/majiayu000/cc-starship/pkg/utils"
)

// HealthResponse represents a health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// HealthCheck handles the health check request
func HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   "1.0.0", // This should be loaded from a config or build info
	}

	utils.JSONResponse(c, http.StatusOK, response)
}
