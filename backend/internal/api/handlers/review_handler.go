package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/majiayu000/cc-starship/internal/core/domain"
	"github.com/majiayu000/cc-starship/internal/core/ports"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

// ReviewHandler handles HTTP requests for review operations
type ReviewHandler struct {
	service ports.ReviewService
	logger  *logger.Logger
}

// NewReviewHandler creates a new ReviewHandler
func NewReviewHandler(service ports.ReviewService, logger *logger.Logger) *ReviewHandler {
	return &ReviewHandler{
		service: service,
		logger:  logger,
	}
}

// GetByID handles GET requests to retrieve a specific reviewable item
func (h *ReviewHandler) GetByID(c *gin.Context) {
	// Get the data source and ID from the URI
	dataSource := c.Param("dataSource")
	originalID := c.Param("id")

	// Validate data source
	if !domain.IsValidDataSource(dataSource) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid data source: %s", dataSource),
		})
		return
	}

	// Call the service
	item, err := h.service.GetByOriginalID(c.Request.Context(), dataSource, originalID)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error retrieving item: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve item",
		})
		return
	}

	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Item with ID %s not found", originalID),
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

// GetAll handles GET requests to retrieve all reviewable items with filtering and pagination
func (h *ReviewHandler) GetAll(c *gin.Context) {
	// Get the data source from the URI
	dataSource := c.Param("dataSource")

	// Validate data source
	if !domain.IsValidDataSource(dataSource) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid data source: %s", dataSource),
		})
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	reviewStatus := c.Query("reviewStatus")
	domainFilter := c.Query("domain")
	skill := c.Query("skill")
	searchTerm := c.Query("searchTerm")
	sortBy := c.Query("sortBy")
	sortOrder := c.Query("sortOrder")

	// Create query params
	params := domain.QueryParams{
		DataSource:   dataSource,
		Page:         page,
		PageSize:     pageSize,
		ReviewStatus: reviewStatus,
		Domain:       domainFilter,
		Skill:        skill,
		SearchTerm:   searchTerm,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
	}

	// Call the service
	result, err := h.service.GetAll(c.Request.Context(), params)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error retrieving items: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve items",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ReviewItem handles POST requests to update the review status of an item
func (h *ReviewHandler) ReviewItem(c *gin.Context) {
	// 获取 URI 参数
	dataSource := c.Param("dataSource")
	originalID := c.Param("id")

	// 记录请求信息
	h.logger.Info(fmt.Sprintf("Received review request for dataSource=%s, originalID=%s",
		dataSource, originalID))

	// 验证数据源
	if !domain.IsValidDataSource(dataSource) {
		h.logger.Warn(fmt.Sprintf("Invalid data source: %s", dataSource))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid data source: %s", dataSource),
		})
		return
	}

	// 解析请求体
	var update domain.ReviewStatusUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		h.logger.Warn(fmt.Sprintf("Invalid request body: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// 记录审核信息
	h.logger.Info(fmt.Sprintf("Updating item %s to status: %s", originalID, update.Status))

	// 验证审核状态
	if !domain.IsValidReviewStatus(update.Status) {
		h.logger.Warn(fmt.Sprintf("Invalid review status: %s", update.Status))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid review status: %s", update.Status),
		})
		return
	}

	// 调用服务
	err := h.service.ReviewItem(c.Request.Context(), dataSource, originalID, update)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error updating review status: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update review status",
			"details": err.Error(), // 添加详细错误信息
		})
		return
	}

	h.logger.Info(fmt.Sprintf("Successfully updated item %s to status %s", originalID, update.Status))
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Review status of item %s updated to %s", originalID, update.Status),
	})
}

// GetFilterOptions handles GET requests to retrieve filter options
func (h *ReviewHandler) GetFilterOptions(c *gin.Context) {
	// Get the data source from the URI
	dataSource := c.Param("dataSource")

	// Validate data source
	if !domain.IsValidDataSource(dataSource) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid data source: %s", dataSource),
		})
		return
	}

	// Call the service
	options, err := h.service.GetFilterOptions(c.Request.Context(), dataSource)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error retrieving filter options: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve filter options",
		})
		return
	}

	c.JSON(http.StatusOK, options)
}

// GetDataSources handles GET requests to retrieve available data sources
func (h *ReviewHandler) GetDataSources(c *gin.Context) {
	// Call the service
	sources, err := h.service.GetDataSources(c.Request.Context())
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error retrieving data sources: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve data sources",
		})
		return
	}

	c.JSON(http.StatusOK, sources)
}
