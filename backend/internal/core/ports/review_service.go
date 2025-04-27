package ports

import (
	"context"

	"github.com/majiayu000/cc-starship/internal/core/domain"
)

// GenericService defines a generic service interface for any entity
type GenericService interface {
	// GetByID retrieves an item by its ID
	GetByID(ctx context.Context, id string) (*domain.GenericItem, error)

	// GetAll retrieves all items with pagination and filtering
	GetAll(ctx context.Context, params domain.QueryParams) (*domain.PaginatedResult, error)

	// Create creates a new item
	Create(ctx context.Context, item *domain.GenericItem) error

	// Update updates an existing item
	Update(ctx context.Context, item *domain.GenericItem) error

	// UpdateStatus updates the status of an item
	UpdateStatus(ctx context.Context, id string, update domain.StatusUpdate) error

	// Delete deletes an item
	Delete(ctx context.Context, id string) error
}
