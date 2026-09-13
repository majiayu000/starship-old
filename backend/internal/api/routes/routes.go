package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/majiayu000/cc-starship/internal/api/handlers"
	"github.com/majiayu000/cc-starship/internal/api/middleware"
	"github.com/majiayu000/cc-starship/internal/core/ports"
	"github.com/majiayu000/cc-starship/internal/infrastructure/cache"
	"github.com/majiayu000/cc-starship/pkg/config"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(
	router *gin.Engine,
	userService ports.UserService,
	authService ports.AuthService,
	cacheService *cache.CacheService,
	reviewService ports.ReviewService,
	cfg *config.Config,
	logger *logger.Logger,
) {
	// Create handlers
	var userHandler *handlers.UserHandler
	var authHandler *handlers.AuthHandler
	var cacheHandler *handlers.CacheHandler
	var reviewHandler *handlers.ReviewHandler

	// Create middlewares
	var authMiddleware gin.HandlerFunc
	var adminRoleMiddleware gin.HandlerFunc

	// Initialize JWT auth handlers/middleware only when auth is enabled and services exist.
	// When auth.enabled=false, do not register /auth or JWT middleware (avoids empty-key tokens).
	if cfg.Auth.Enabled && userService != nil && authService != nil {
		userHandler = handlers.NewUserHandler(userService, logger)
		authHandler = handlers.NewAuthHandler(authService, logger)
		authMiddleware = middleware.Auth(authService)
		adminRoleMiddleware = middleware.RequireRole("admin")
	}

	// Initialize cache handler if caching is enabled
	if cacheService != nil && cfg.Features.EnableCaching {
		cacheHandler = handlers.NewCacheHandler(cacheService, logger)
	}

	// Initialize review handler if review service is available
	if reviewService != nil {
		reviewHandler = handlers.NewReviewHandler(reviewService, logger)
	}

	// Register health check route
	router.GET("/health", handlers.HealthCheck)

	// Group API routes
	api := router.Group("/api/v1")
	{
		// Register auth routes only when auth is enabled
		if cfg.Auth.Enabled && authHandler != nil {
			auth := api.Group("/auth")
			{
				auth.POST("/register", authHandler.Register)
				auth.POST("/login", authHandler.Login)
			}
		}

		// Register user routes only when auth is enabled (JWT-protected)
		if cfg.Auth.Enabled && userHandler != nil {
			users := api.Group("/users")
			users.Use(authMiddleware)
			{
				users.GET("", adminRoleMiddleware, userHandler.GetUsers)
				users.GET("/:id", userHandler.GetUser)
				users.POST("", adminRoleMiddleware, userHandler.CreateUser)
				users.PUT("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", adminRoleMiddleware, userHandler.DeleteUser)
			}
		}

		// Register cache routes if cache service is available
		if cacheHandler != nil {
			cache := api.Group("/cache")
			if cfg.Auth.Enabled && authService != nil {
				cache.Use(authMiddleware)
				cache.Use(adminRoleMiddleware)
			}
			{
				cache.POST("", cacheHandler.Set)
				cache.GET("/:key", cacheHandler.Get)
				cache.DELETE("/:key", cacheHandler.Delete)
			}
		}

		// Register review routes if review service is available
		if reviewHandler != nil {
			// Group for review-related endpoints
			review := api.Group("/review")
			// Apply auth middleware if auth service is available
			// 临时注释掉认证中间件
			// if authService != nil {
			// 	review.Use(authMiddleware)
			// }
			{
				// Get available data sources
				review.GET("/sources", reviewHandler.GetDataSources)

				// Data source specific routes
				dataSource := review.Group("/:dataSource")
				{
					// Get filter options for a data source
					dataSource.GET("/filters", reviewHandler.GetFilterOptions)

					// Get all items with pagination and filtering
					dataSource.GET("/items", reviewHandler.GetAll)

					// Get a specific item
					dataSource.GET("/items/:id", reviewHandler.GetByID)

					// Review an item (update status)
					dataSource.POST("/items/:id/review", reviewHandler.ReviewItem)
				}
			}
		}
	}
}
