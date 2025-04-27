package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/majiayu000/cc-starship/internal/api/middleware"
	"github.com/majiayu000/cc-starship/internal/api/routes"
	"github.com/majiayu000/cc-starship/internal/core/ports"
	"github.com/majiayu000/cc-starship/internal/infrastructure/cache"
	"github.com/majiayu000/cc-starship/pkg/config"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

// Server represents the HTTP server
type Server struct {
	router         *gin.Engine
	server         *http.Server
	cfg            *config.Config
	logger         *logger.Logger
	logManager     *logger.LogManager
	userService    ports.UserService
	authService    ports.AuthService
	cacheService   *cache.CacheService
	genericService ports.GenericService
}

// NewServer creates a new HTTP server
func NewServer(
	cfg *config.Config,
	logger *logger.Logger,
	userService ports.UserService,
	authService ports.AuthService,
	logManager *logger.LogManager,
	cacheService *cache.CacheService,
	genericService ports.GenericService,
) *Server {
	// Set Gin mode based on environment
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else if cfg.App.Environment == "test" {
		gin.SetMode(gin.TestMode)
	}

	// Create a new Gin router
	router := gin.New()

	// Add recovery middleware to handle panics
	router.Use(gin.Recovery())

	// 添加日志中间件
	router.Use(middleware.Logging(logManager))

	// Apply CORS middleware if enabled
	if cfg.Features.EnableCORS {
		router.Use(middleware.CORS(nil)) // Use default CORS config
	}

	// Create server
	server := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &Server{
		router:         router,
		server:         server,
		cfg:            cfg,
		logger:         logger,
		logManager:     logManager,
		userService:    userService,
		authService:    authService,
		cacheService:   cacheService,
		genericService: genericService,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Register routes
	routes.RegisterRoutes(
		s.router,
		s.userService,
		s.authService,
		s.cacheService,
		s.cfg,
		s.logger,
	)

	// Register generic routes
	routes.RegisterGenericRoutes(
		s.router,
		s.genericService,
		s.authService,
		s.cfg,
		s.logger,
	)

	// Log server start
	s.logger.Info("Starting HTTP server on " + s.cfg.Server.Address)

	// Start server
	return s.server.ListenAndServe()
}

// Stop stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server")
	return s.server.Shutdown(ctx)
}
