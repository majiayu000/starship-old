package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/majiayu000/cc-starship/internal/infrastructure/cache"
	"github.com/majiayu000/cc-starship/pkg/errors"
	"github.com/majiayu000/cc-starship/pkg/logger"
	"github.com/majiayu000/cc-starship/pkg/utils"
)

// CacheHandler handles cache-related requests
type CacheHandler struct {
	cacheService *cache.CacheService
	logger       *logger.Logger
}

// NewCacheHandler creates a new CacheHandler
func NewCacheHandler(cacheService *cache.CacheService, logger *logger.Logger) *CacheHandler {
	return &CacheHandler{
		cacheService: cacheService,
		logger:       logger,
	}
}

// CacheRequest represents a request to set a cache value
type CacheRequest struct {
	Key   string      `json:"key" binding:"required"`
	Value interface{} `json:"value" binding:"required"`
	TTL   string      `json:"ttl" binding:"required"`
}

// CacheResponse represents a cache response
type CacheResponse struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value,omitempty"`
	Found bool        `json:"found"`
}

// Set handles setting a value in the cache
func (h *CacheHandler) Set(c *gin.Context) {
	// Parse request body
	var req CacheRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequest("Invalid request body", err))
		return
	}

	// Parse TTL
	ttl, err := time.ParseDuration(req.TTL)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequest("Invalid TTL format", err))
		return
	}

	// Set value in cache
	if err := h.cacheService.Set(c, req.Key, req.Value, ttl); err != nil {
		utils.ErrorResponse(c, errors.NewInternal("Failed to set cache value", err))
		return
	}

	// Send response
	utils.JSONResponse(c, http.StatusOK, map[string]string{
		"message": "Value cached successfully",
	})
}

// Get handles getting a value from the cache
func (h *CacheHandler) Get(c *gin.Context) {
	// Get key from request path parameters
	key := c.Param("key")
	if key == "" {
		utils.ErrorResponse(c, errors.NewBadRequest("Key is required", nil))
		return
	}

	// Get value from cache
	var value interface{}
	found, err := h.cacheService.Get(c, key, &value)
	if err != nil {
		utils.ErrorResponse(c, errors.NewInternal("Failed to get cache value", err))
		return
	}

	// Send response
	response := CacheResponse{
		Key:   key,
		Found: found,
	}

	if found {
		response.Value = value
	}

	utils.JSONResponse(c, http.StatusOK, response)
}

// Delete handles deleting a value from the cache
func (h *CacheHandler) Delete(c *gin.Context) {
	// Get key from request path parameters
	key := c.Param("key")
	if key == "" {
		utils.ErrorResponse(c, errors.NewBadRequest("Key is required", nil))
		return
	}

	// Delete value from cache
	if err := h.cacheService.Delete(c, key); err != nil {
		utils.ErrorResponse(c, errors.NewInternal("Failed to delete cache value", err))
		return
	}

	// Send response
	utils.JSONResponse(c, http.StatusOK, map[string]string{
		"message": "Value deleted successfully",
	})
}
