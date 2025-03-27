package services

import (
	"context"
	"fmt"

	"github.com/majiayu000/cc-starship/internal/core/domain"
	"github.com/majiayu000/cc-starship/internal/core/ports"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

// ReviewService implements the service for handling reviewable items
type ReviewService struct {
	repository ports.ReviewRepository
	logger     *logger.Logger
}

// NewReviewService creates a new ReviewService
func NewReviewService(repository ports.ReviewRepository, logger *logger.Logger) *ReviewService {
	return &ReviewService{
		repository: repository,
		logger:     logger,
	}
}

// GetByOriginalID retrieves a reviewable item by its original ID and data source
func (s *ReviewService) GetByOriginalID(ctx context.Context, dataSource, originalID string) (interface{}, error) {
	// Validate data source
	if !domain.IsValidDataSource(dataSource) {
		return nil, fmt.Errorf("invalid data source: %s", dataSource)
	}

	// Log the operation
	s.logger.Info(fmt.Sprintf("Retrieving item with original ID %s from data source %s", originalID, dataSource))

	// Call the repository
	return s.repository.FindByOriginalID(ctx, dataSource, originalID)
}

// GetAll retrieves all reviewable items with pagination and filtering
func (s *ReviewService) GetAll(ctx context.Context, params domain.QueryParams) (*domain.PaginatedResult, error) {
	// Validate data source
	if !domain.IsValidDataSource(params.DataSource) {
		return nil, fmt.Errorf("invalid data source: %s", params.DataSource)
	}

	// Log the operation
	s.logger.Info(fmt.Sprintf("Retrieving items from data source %s with params: %+v", params.DataSource, params))

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

// ReviewItem reviews an item
func (s *ReviewService) ReviewItem(ctx context.Context, dataSource, originalID string, update domain.ReviewStatusUpdate) error {
	// Validate data source
	if !domain.IsValidDataSource(dataSource) {
		return fmt.Errorf("invalid data source: %s", dataSource)
	}

	// Validate review status
	if !domain.IsValidReviewStatus(update.Status) {
		return fmt.Errorf("invalid review status: %s", update.Status)
	}

	// Log the operation
	s.logger.Info(fmt.Sprintf("Updating review status of item with original ID %s from data source %s to %s",
		originalID, dataSource, update.Status))

	// Call the repository
	return s.repository.UpdateReviewStatus(ctx, dataSource, originalID, update)
}

// GetFilterOptions retrieves available options for filtering
func (s *ReviewService) GetFilterOptions(ctx context.Context, dataSource string) (map[string][]string, error) {
	// Validate data source
	if !domain.IsValidDataSource(dataSource) {
		return nil, fmt.Errorf("invalid data source: %s", dataSource)
	}

	// Log the operation
	s.logger.Info(fmt.Sprintf("Retrieving filter options for data source %s", dataSource))

	// Call the repository
	return s.repository.GetFilterOptions(ctx, dataSource)
}

// GetDataSources returns a list of available data sources
func (s *ReviewService) GetDataSources(ctx context.Context) ([]string, error) {
	// Log the operation
	s.logger.Info("Retrieving available data sources")

	// Call the repository
	return s.repository.GetDataSources(ctx)
}
