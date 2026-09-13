package services

import (
	"context"
	"net/mail"
	"strings"

	"github.com/majiayu000/cc-starship/internal/core/domain"
	"github.com/majiayu000/cc-starship/internal/core/ports"
	"github.com/majiayu000/cc-starship/internal/infrastructure/auth"
	"github.com/majiayu000/cc-starship/pkg/config"
	"github.com/majiayu000/cc-starship/pkg/errors"
	"github.com/majiayu000/cc-starship/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

// AuthService implements the authentication service
type AuthService struct {
	userRepo   ports.UserRepository
	jwtService *auth.JWTService
	logger     *logger.Logger
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo ports.UserRepository, jwtService *auth.JWTService, logger *logger.Logger) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
		logger:     logger,
	}
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, email, password, firstName, lastName string) (*domain.User, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.NewConflict("User with this email already exists", nil)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewInternal("Failed to hash password", err)
	}

	// Create new user (always role=user). First admin is provisioned only via
	// EnsureBootstrapAdmin using deployment-controlled credentials.
	user := domain.NewUser(email, string(hashedPassword), firstName, lastName)

	// Save user to repository
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.NewInternal("Failed to create user", err)
	}

	return user, nil
}

// EnsureBootstrapAdmin creates the first admin from deployment config when no
// admin exists. Safe to call on every startup; concurrent callers serialize via
// CreateAdminIfAbsent. Public registration never promotes to admin.
func (s *AuthService) EnsureBootstrapAdmin(ctx context.Context, cfg config.BootstrapAdminConfig) error {
	email := strings.TrimSpace(cfg.Email)
	password := cfg.Password
	emailEmpty := email == ""
	passwordEmpty := password == ""
	// Only skip when both credentials are absent. A partial config is a startup
	// error so clean deployments cannot silently start without an admin.
	if emailEmpty && passwordEmpty {
		return nil
	}
	if emailEmpty || passwordEmpty {
		return errors.NewBadRequest("Bootstrap admin requires both email and password", nil)
	}

	// Require a plain login email matching LoginRequest binding:"email".
	// mail.ParseAddress accepts display-name forms like `Admin <a@b.com>` which
	// cannot be typed into the HTML email input or pass gin email validation.
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Name != "" || !strings.EqualFold(parsed.Address, email) {
		return errors.NewBadRequest("Invalid bootstrap admin email address", err)
	}
	email = parsed.Address

	firstName := strings.TrimSpace(cfg.FirstName)
	if firstName == "" {
		firstName = "Admin"
	}
	lastName := strings.TrimSpace(cfg.LastName)
	if lastName == "" {
		lastName = "User"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewInternal("Failed to hash bootstrap admin password", err)
	}

	user := domain.NewUser(email, string(hashedPassword), firstName, lastName)
	user.Role = "admin"

	created, err := s.userRepo.CreateAdminIfAbsent(ctx, user)
	if err != nil {
		return errors.NewInternal("Failed to bootstrap admin user", err)
	}
	if created {
		s.logger.Info("Bootstrapped first admin user from deployment config: " + email)
	}
	return nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.NewUnauthorized("Invalid email or password", nil)
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.NewUnauthorized("Invalid email or password", nil)
	}

	if !user.Active {
		return "", errors.NewUnauthorized("Account is inactive", nil)
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", errors.NewInternal("Failed to generate token", err)
	}

	return token, nil
}

// ValidateToken validates a JWT token and returns the associated user
func (s *AuthService) ValidateToken(ctx context.Context, token string) (*domain.User, error) {
	// Validate token
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, errors.NewUnauthorized("Invalid token", err)
	}

	// Get user by ID
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.NewUnauthorized("User not found", err)
	}

	if !user.Active {
		return nil, errors.NewUnauthorized("Account is inactive", nil)
	}

	return user, nil
}
