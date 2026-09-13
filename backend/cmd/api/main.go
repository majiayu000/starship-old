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
	var reviewService ports.ReviewService

	// Initialize user and auth services only if PostgreSQL is available
	if dbManager.PostgreSQL != nil {
		// Create repositories
		userRepo := postgres.NewUserRepository(dbManager.PostgreSQL, l)

		// Create review repository
		reviewRepo := postgres.NewReviewRepository(dbManager.PostgreSQL, l)

		// Create review service
		reviewService = services.NewReviewService(reviewRepo, l)

		// Create services
		userService = services.NewUserService(userRepo, l)

		// Honor auth.enabled and features.enable_auth so deployments can run
		// PostgreSQL-backed review without forcing Bearer auth on every route.
		authEnabled := cfg.Auth.Enabled && cfg.Features.EnableAuth
		if authEnabled {
			jwtService := auth.NewJWTService(cfg)
			authSvc := services.NewAuthService(userRepo, jwtService, l)
			if err := authSvc.EnsureBootstrapAdmin(context.Background(), cfg.Auth.BootstrapAdmin); err != nil {
				l.Fatal("Failed to ensure bootstrap admin: " + err.Error())
			}
			authService = authSvc

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
			l.Info("Authentication disabled by configuration (auth.enabled / features.enable_auth)")

			l.LogBusiness(
				logger.LogLevelInfo,
				"system_init",
				"system",
				"auth",
				map[string]interface{}{
					"service": "PostgreSQL",
					"status":  "disabled",
					"reason":  "auth.enabled or features.enable_auth is false",
				},
			)
		}
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
	server := api.NewServer(cfg, l, userService, authService, logManager, cacheService, reviewService)

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
