package services

import (
	"context"
	stderrors "errors"

	"github.com/majiayu000/cc-starship/internal/core/domain"
	"github.com/majiayu000/cc-starship/internal/core/ports"
	"github.com/majiayu000/cc-starship/pkg/errors"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

// UserService implements the user service
type UserService struct {
	userRepo ports.UserRepository
	logger   *logger.Logger
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo ports.UserRepository, logger *logger.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// GetUserByID gets a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFound("User not found", err)
	}
	return user, nil
}

// GetUserByEmail gets a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.NewNotFound("User not found", err)
	}
	return user, nil
}

// GetUsers gets all users
func (s *UserService) GetUsers(ctx context.Context) ([]*domain.User, error) {
	return s.userRepo.FindAll(ctx)
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	// Check if user with the same email already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return errors.NewConflict("User with this email already exists", nil)
	}

	// Create the user
	return s.userRepo.Create(ctx, user)
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, user *domain.User) error {
	// Check if user exists
	_, err := s.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		return errors.NewNotFound("User not found", err)
	}

	// Update the user (repository enforces last-active-admin invariant atomically)
	if err := s.userRepo.Update(ctx, user); err != nil {
		if stderrors.Is(err, domain.ErrCannotDemoteLastAdmin) {
			return errors.NewForbidden("Cannot demote or deactivate the sole active administrator", err)
		}
		return err
	}
	return nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	// Check if user exists
	_, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return errors.NewNotFound("User not found", err)
	}

	// Delete the user
	return s.userRepo.Delete(ctx, id)
}
