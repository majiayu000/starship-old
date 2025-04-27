package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/majiayu000/cc-starship/internal/core/domain"
	"github.com/majiayu000/cc-starship/internal/core/ports"
	"github.com/majiayu000/cc-starship/pkg/logger"
	"github.com/majiayu000/cc-starship/pkg/utils"
)

// GenericHandler handles generic item requests
type GenericHandler struct {
	service ports.GenericService
	logger  *logger.Logger
}

// NewGenericHandler creates a new GenericHandler
func NewGenericHandler(service ports.GenericService, logger *logger.Logger) *GenericHandler {
	return &GenericHandler{
		service: service,
		logger:  logger,
	}
}

// GetByID handles the request to get an item by ID
func (h *GenericHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	// Log the request
	h.logger.Info(fmt.Sprintf("Getting item with ID: %s", id))

	// Call the service
	item, err := h.service.GetByID(c, id)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("Error getting item: %v", err))
		utils.JSONResponse(c, http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("item not found: %v", err),
		})
		return
	}

	// Return the item
	utils.JSONResponse(c, http.StatusOK, item)
}

// GetAll handles the request to get all items with pagination and filtering
func (h *GenericHandler) GetAll(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")
	searchTerm := c.Query("search")
	sortBy := c.DefaultQuery("sortBy", "createdAt")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	// Log the request
	h.logger.Info(fmt.Sprintf("Getting items with page: %d, pageSize: %d, status: %s, search: %s, sortBy: %s, sortOrder: %s",
		page, pageSize, status, searchTerm, sortBy, sortOrder))

	// Create query params
	params := domain.QueryParams{
		Page:       page,
		PageSize:   pageSize,
		Status:     status,
		SearchTerm: searchTerm,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	}

	// Call the service
	result, err := h.service.GetAll(c, params)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("Error getting items: %v", err))
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error getting items: %v", err),
		})
		return
	}

	// Return the result
	utils.JSONResponse(c, http.StatusOK, result)
}

// Create handles the request to create a new item
func (h *GenericHandler) Create(c *gin.Context) {
	// Parse request body
	var item domain.GenericItem
	if err := c.ShouldBindJSON(&item); err != nil {
		h.logger.Warn(fmt.Sprintf("Invalid request body: %v", err))
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request body: %v", err),
		})
		return
	}

	// Log the request
	h.logger.Info(fmt.Sprintf("Creating item with name: %s", item.Name))

	// Call the service
	err := h.service.Create(c, &item)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("Error creating item: %v", err))
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error creating item: %v", err),
		})
		return
	}

	// Return success
	utils.JSONResponse(c, http.StatusCreated, item)
}

// Update handles the request to update an existing item
func (h *GenericHandler) Update(c *gin.Context) {
	id := c.Param("id")

	// Parse request body
	var item domain.GenericItem
	if err := c.ShouldBindJSON(&item); err != nil {
		h.logger.Warn(fmt.Sprintf("Invalid request body: %v", err))
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request body: %v", err),
		})
		return
	}

	// Ensure ID in path matches ID in body
	item.ID = id

	// Log the request
	h.logger.Info(fmt.Sprintf("Updating item with ID: %s", id))

	// Call the service
	err := h.service.Update(c, &item)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("Error updating item: %v", err))
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error updating item: %v", err),
		})
		return
	}

	// Return success
	utils.JSONResponse(c, http.StatusOK, item)
}

// UpdateStatus handles the request to update an item's status
func (h *GenericHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	// Parse request body
	var update domain.StatusUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		h.logger.Warn(fmt.Sprintf("Invalid request body: %v", err))
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request body: %v", err),
		})
		return
	}

	// Log the request
	h.logger.Info(fmt.Sprintf("Updating status of item with ID: %s to %s", id, update.Status))

	// Validate status
	if !domain.IsValidStatus(update.Status) {
		h.logger.Warn(fmt.Sprintf("Invalid status: %s", update.Status))
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid status: %s", update.Status),
		})
		return
	}

	// Call the service
	err := h.service.UpdateStatus(c, id, update)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("Error updating item status: %v", err))
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error updating item status: %v", err),
		})
		return
	}

	// Return success
	utils.JSONResponse(c, http.StatusOK, gin.H{
		"message": fmt.Sprintf("Status of item %s updated to %s", id, update.Status),
	})
}

// Delete handles the request to delete an item
func (h *GenericHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	// Log the request
	h.logger.Info(fmt.Sprintf("Deleting item with ID: %s", id))

	// Call the service
	err := h.service.Delete(c, id)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("Error deleting item: %v", err))
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error deleting item: %v", err),
		})
		return
	}

	// Return success
	utils.JSONResponse(c, http.StatusOK, gin.H{
		"message": fmt.Sprintf("Item %s deleted successfully", id),
	})
}
