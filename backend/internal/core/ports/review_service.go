package ports

import (
	"context"

	"github.com/majiayu000/cc-starship/internal/core/domain"
)

// ReviewService defines the interface for review business logic
type ReviewService interface {
	// GetByOriginalID retrieves a reviewable item by its original ID and data source
	GetByOriginalID(ctx context.Context, dataSource, originalID string) (interface{}, error)

	// GetAll retrieves all reviewable items with pagination and filtering
	GetAll(ctx context.Context, params domain.QueryParams) (*domain.PaginatedResult, error)

	// ReviewItem reviews an item
	ReviewItem(ctx context.Context, dataSource, originalID string, update domain.ReviewStatusUpdate) error

	// GetFilterOptions retrieves available options for filtering (domains, skills, etc.)
	GetFilterOptions(ctx context.Context, dataSource string) (map[string][]string, error)

	// GetDataSources returns a list of available data sources
	GetDataSources(ctx context.Context) ([]string, error)
}
