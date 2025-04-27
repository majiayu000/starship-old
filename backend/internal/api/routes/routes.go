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
	cfg *config.Config,
	logger *logger.Logger,
) {
	// Create handlers
	var userHandler *handlers.UserHandler
	var authHandler *handlers.AuthHandler
	var cacheHandler *handlers.CacheHandler

	// Create middlewares
	var authMiddleware gin.HandlerFunc
	var adminRoleMiddleware gin.HandlerFunc

	// Initialize handlers and middlewares only if services are available
	if userService != nil && authService != nil {
		userHandler = handlers.NewUserHandler(userService, logger)
		authHandler = handlers.NewAuthHandler(authService, logger)
		authMiddleware = middleware.Auth(authService)
		adminRoleMiddleware = middleware.RequireRole("admin")
	}

	// Initialize cache handler if caching is enabled
	if cacheService != nil && cfg.Features.EnableCaching {
		cacheHandler = handlers.NewCacheHandler(cacheService, logger)
	}

	// Register health check route
	router.GET("/health", handlers.HealthCheck)

	// Group API routes
	api := router.Group("/api/v1")
	{
		// Register auth routes if auth service is available
		if authHandler != nil {
			auth := api.Group("/auth")
			{
				auth.POST("/register", authHandler.Register)
				auth.POST("/login", authHandler.Login)
			}
		}

		// Register user routes if user service is available
		if userHandler != nil {
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
			if authService != nil {
				cache.Use(authMiddleware)
				cache.Use(adminRoleMiddleware)
			}
			{
				cache.POST("", cacheHandler.Set)
				cache.GET("/:key", cacheHandler.Get)
				cache.DELETE("/:key", cacheHandler.Delete)
			}
		}

		// Register generic routes here if needed
	}
}
