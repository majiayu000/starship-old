package services

import (
	"context"
	"fmt"

	"github.com/majiayu000/cc-starship/internal/core/domain"
	"github.com/majiayu000/cc-starship/internal/core/ports"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

// GenericService implements the service for handling generic items
type GenericService struct {
	repository ports.GenericRepository
	logger     *logger.Logger
}

// NewGenericService creates a new GenericService
func NewGenericService(repository ports.GenericRepository, logger *logger.Logger) *GenericService {
	return &GenericService{
		repository: repository,
		logger:     logger,
	}
}

// GetByID retrieves an item by its ID
func (s *GenericService) GetByID(ctx context.Context, id string) (*domain.GenericItem, error) {
	// Log the operation
	s.logger.Info(fmt.Sprintf("Retrieving item with ID %s", id))

	// Call the repository
	return s.repository.FindByID(ctx, id)
}

// GetAll retrieves all items with pagination and filtering
func (s *GenericService) GetAll(ctx context.Context, params domain.QueryParams) (*domain.PaginatedResult, error) {
	// Log the operation
	s.logger.Info(fmt.Sprintf("Retrieving items with params: %+v", params))

	// Set default pagination values if not provided
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	// Call the repository
	return s.repository.FindAll(ctx, params)
}

// Create creates a new item
func (s *GenericService) Create(ctx context.Context, item *domain.GenericItem) error {
	// Log the operation
	s.logger.Info(fmt.Sprintf("Creating new item with name: %s", item.Name))

	// Call the repository
	return s.repository.Create(ctx, item)
}

// Update updates an existing item
func (s *GenericService) Update(ctx context.Context, item *domain.GenericItem) error {
	// Log the operation
	s.logger.Info(fmt.Sprintf("Updating item with ID: %s", item.ID))

	// Call the repository
	return s.repository.Update(ctx, item)
}

// UpdateStatus updates the status of an item
func (s *GenericService) UpdateStatus(ctx context.Context, id string, update domain.StatusUpdate) error {
	// Validate status
	if !domain.IsValidStatus(update.Status) {
		return fmt.Errorf("invalid status: %s", update.Status)
	}

	// Log the operation
	s.logger.Info(fmt.Sprintf("Updating status of item with ID %s to %s", id, update.Status))

	// Call the repository
	return s.repository.UpdateStatus(ctx, id, update)
}

// Delete deletes an item
func (s *GenericService) Delete(ctx context.Context, id string) error {
	// Log the operation
	s.logger.Info(fmt.Sprintf("Deleting item with ID %s", id))

	// Call the repository
	return s.repository.Delete(ctx, id)
}
