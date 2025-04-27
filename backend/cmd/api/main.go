package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/majiayu000/cc-starship/internal/api"
	"github.com/majiayu000/cc-starship/internal/core/ports"
	"github.com/majiayu000/cc-starship/internal/core/services"
	"github.com/majiayu000/cc-starship/internal/infrastructure/auth"
	"github.com/majiayu000/cc-starship/internal/infrastructure/cache"
	"github.com/majiayu000/cc-starship/internal/infrastructure/database"
	"github.com/majiayu000/cc-starship/internal/repositories/postgres"
	"github.com/majiayu000/cc-starship/pkg/config"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger manager
	logManager, err := logger.NewLogManager(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize log manager: %v", err)
	}
	logManager.Start()
	defer logManager.Stop()

	// Initialize legacy logger for backward compatibility
	l := logger.New(cfg.Logger)
	l.SetLogManager(logManager)

	// Log application start
	l.Info("Starting application: " + cfg.App.Name)

	// Initialize databases
	dbManager, err := database.NewDataAccessManager(cfg, l)
	if err != nil {
		l.Fatal("Failed to initialize database manager: " + err.Error())
	}
	defer dbManager.Close()

	// Initialize cache service
	cacheService := cache.NewCacheService(dbManager, l)

	// Log cache service initialization
	if cfg.Features.EnableCaching {
		l.LogBusiness(
			logger.LogLevelInfo,
			"system_init",
			"system",
			"cache",
			map[string]interface{}{
				"service":     "CacheService",
				"status":      "initialized",
				"using_redis": dbManager.Redis != nil,
			},
		)
	}

	// Initialize services based on available databases
	var userService ports.UserService
	var authService ports.AuthService
	var genericService ports.GenericService

	// Initialize user and auth services only if PostgreSQL is available
	if dbManager.PostgreSQL != nil {
		// Create JWT service
		jwtService := auth.NewJWTService(cfg)

		// Create repositories
		userRepo := postgres.NewUserRepository(dbManager.PostgreSQL, l)

		// Create generic repository
		genericRepo := postgres.NewGenericRepository(dbManager.PostgreSQL, l, "generic_items")

		// Create generic service
		genericService = services.NewGenericService(genericRepo, l)

		// Create services
		userService = services.NewUserService(userRepo, l)
		authService = services.NewAuthService(userRepo, jwtService, l)

		l.LogBusiness(
			logger.LogLevelInfo,
			"system_init",
			"system",
			"auth",
			map[string]interface{}{
				"service": "PostgreSQL",
				"status":  "initialized",
			},
		)
	} else {
		l.Warn("PostgreSQL is not initialized. User authentication features will be disabled.")

		l.LogBusiness(
			logger.LogLevelWarn,
			"system_init",
			"system",
			"auth",
			map[string]interface{}{
				"service": "PostgreSQL",
				"status":  "disabled",
				"reason":  "Database not initialized",
			},
		)
	}

	// Create and start HTTP server
	server := api.NewServer(cfg, l, userService, authService, logManager, cacheService, genericService)

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			l.Fatal("Failed to start server: " + err.Error())
		}
	}()

	// Log server started
	l.LogBusiness(
		logger.LogLevelInfo,
		"server_start",
		"system",
		"http_server",
		map[string]interface{}{
			"address": cfg.Server.Address,
		},
	)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	l.Info("Shutting down server...")

	// Create a deadline to wait for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := server.Stop(ctx); err != nil {
		l.Fatal("Server forced to shutdown: " + err.Error())
	}

	l.Info("Server stopped")

	l.LogBusiness(
		logger.LogLevelInfo,
		"server_stop",
		"system",
		"http_server",
		map[string]interface{}{
			"reason": "Graceful shutdown",
		},
	)
}
