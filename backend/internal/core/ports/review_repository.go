package ports

import (
	"context"

	"github.com/majiayu000/cc-starship/internal/core/domain"
)

// ReviewRepository defines the interface for interacting with reviewable items
type ReviewRepository interface {
	// FindByOriginalID retrieves a reviewable item by its original ID and data source
	FindByOriginalID(ctx context.Context, dataSource, originalID string) (interface{}, error)

	// FindAll retrieves all reviewable items with pagination and filtering
	FindAll(ctx context.Context, params domain.QueryParams) (*domain.PaginatedResult, error)

	// UpdateReviewStatus updates the review status of a reviewable item
	UpdateReviewStatus(ctx context.Context, dataSource, originalID string, update domain.ReviewStatusUpdate) error

	// GetFilterOptions retrieves available options for filtering (domains, skills, etc.)
	GetFilterOptions(ctx context.Context, dataSource string) (map[string][]string, error)

	// GetDataSources returns a list of available data sources
	GetDataSources(ctx context.Context) ([]string, error)
}
