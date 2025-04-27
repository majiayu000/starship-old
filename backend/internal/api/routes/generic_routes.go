package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/majiayu000/cc-starship/internal/api/handlers"
	"github.com/majiayu000/cc-starship/internal/api/middleware"
	"github.com/majiayu000/cc-starship/internal/core/ports"
	"github.com/majiayu000/cc-starship/pkg/config"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

// RegisterGenericRoutes registers routes for generic items
func RegisterGenericRoutes(
	router *gin.Engine,
	genericService ports.GenericService,
	authService ports.AuthService,
	cfg *config.Config,
	logger *logger.Logger,
) {
	// Create handler
	genericHandler := handlers.NewGenericHandler(genericService, logger)

	// Create API group
	api := router.Group("/api/v1")

	// Add middleware
	if cfg.Features.EnableAuth {
		api.Use(middleware.Auth(authService))
	}

	// Register routes
	items := api.Group("/items")
	{
		// Get all items with pagination and filtering
		items.GET("", genericHandler.GetAll)

		// Get a specific item
		items.GET("/:id", genericHandler.GetByID)

		// Create a new item
		items.POST("", genericHandler.Create)

		// Update an existing item
		items.PUT("/:id", genericHandler.Update)

		// Update an item's status
		items.PATCH("/:id/status", genericHandler.UpdateStatus)

		// Delete an item
		items.DELETE("/:id", genericHandler.Delete)
	}
}
