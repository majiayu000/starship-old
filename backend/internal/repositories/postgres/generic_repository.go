package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/majiayu000/cc-starship/internal/core/domain"
	"github.com/majiayu000/cc-starship/internal/infrastructure/database"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

// GenericRepository implements the generic repository using PostgreSQL
type GenericRepository struct {
	db        *database.PostgresDB
	logger    *logger.Logger
	tableName string
}

// NewGenericRepository creates a new GenericRepository
func NewGenericRepository(db *database.PostgresDB, logger *logger.Logger, tableName string) *GenericRepository {
	return &GenericRepository{
		db:        db,
		logger:    logger,
		tableName: tableName,
	}
}

// FindByID retrieves an item by its ID
func (r *GenericRepository) FindByID(ctx context.Context, id string) (*domain.GenericItem, error) {
	query := fmt.Sprintf(`SELECT id, name, description, status, created_at, updated_at FROM %s WHERE id = $1`, r.tableName)
	
	row := r.db.QueryRowContext(ctx, query, id)
	
	return r.scanItem(row)
}

// FindAll retrieves all items with pagination and filtering
func (r *GenericRepository) FindAll(ctx context.Context, params domain.QueryParams) (*domain.PaginatedResult, error) {
	// Build the query
	conditions := []string{"1=1"} // Always true condition to simplify query building
	args := []interface{}{}
	paramIdx := 1

	// Add status filter
	if params.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", paramIdx))
		args = append(args, params.Status)
		paramIdx++
	}

	// Add search term filter
	if params.SearchTerm != "" {
		searchCondition := fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", paramIdx, paramIdx)
		conditions = append(conditions, searchCondition)
		args = append(args, "%"+params.SearchTerm+"%")
		paramIdx++
	}

	// Build the WHERE clause
	whereClause := strings.Join(conditions, " AND ")

	// Build the ORDER BY clause
	orderBy := "created_at DESC" // Default sorting
	if params.SortBy != "" {
		direction := "ASC"
		if strings.ToLower(params.SortOrder) == "desc" {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", params.SortBy, direction)
	}

	// Count total items
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, r.tableName, whereClause)
	var totalItems int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("error counting items: %w", err)
	}

	// Calculate pagination
	offset := (params.Page - 1) * params.PageSize
	limit := params.PageSize

	// Add pagination to args
	args = append(args, limit, offset)

	// Build the final query
	query := fmt.Sprintf(`
		SELECT id, name, description, status, created_at, updated_at
		FROM %s
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, r.tableName, whereClause, orderBy, paramIdx, paramIdx+1)

	// Execute the query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying items: %w", err)
	}
	defer rows.Close()

	// Scan the results
	items, err := r.scanItems(rows)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := (totalItems + int64(params.PageSize) - 1) / int64(params.PageSize)

	// Return paginated result
	return &domain.PaginatedResult{
		Items:      items,
		TotalItems: totalItems,
		TotalPages: int(totalPages),
		Page:       params.Page,
		PageSize:   params.PageSize,
	}, nil
}

// Create creates a new item
func (r *GenericRepository) Create(ctx context.Context, item *domain.GenericItem) error {
	// Generate a new ID if not provided
	if item.ID == "" {
		item.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	// Build the query
	query := fmt.Sprintf(`
		INSERT INTO %s (id, name, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, r.tableName)

	// Execute the query
	_, err := r.db.ExecContext(ctx, query,
		item.ID,
		item.Name,
		item.Description,
		item.Status,
		item.CreatedAt,
		item.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating item: %w", err)
	}

	return nil
}

// Update updates an existing item
func (r *GenericRepository) Update(ctx context.Context, item *domain.GenericItem) error {
	// Update timestamp
	item.UpdatedAt = time.Now()

	// Build the query
	query := fmt.Sprintf(`
		UPDATE %s
		SET name = $1, description = $2, status = $3, updated_at = $4
		WHERE id = $5
	`, r.tableName)

	// Execute the query
	result, err := r.db.ExecContext(ctx, query,
		item.Name,
		item.Description,
		item.Status,
		item.UpdatedAt,
		item.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating item: %w", err)
	}

	// Check if the item was found
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item with ID %s not found", item.ID)
	}

	return nil
}

// UpdateStatus updates the status of an item
func (r *GenericRepository) UpdateStatus(ctx context.Context, id string, update domain.StatusUpdate) error {
	// Validate the status
	if !domain.IsValidStatus(update.Status) {
		return fmt.Errorf("invalid status: %s", update.Status)
	}

	// Build the query
	query := fmt.Sprintf(`
		UPDATE %s
		SET status = $1, updated_at = $2
		WHERE id = $3
	`, r.tableName)

	// Execute the query
	result, err := r.db.ExecContext(ctx, query,
		update.Status,
		time.Now(),
		id,
	)

	if err != nil {
		return fmt.Errorf("error updating item status: %w", err)
	}

	// Check if the item was found
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item with ID %s not found", id)
	}

	return nil
}

// Delete deletes an item
func (r *GenericRepository) Delete(ctx context.Context, id string) error {
	// Build the query
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, r.tableName)

	// Execute the query
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting item: %w", err)
	}

	// Check if the item was found
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item with ID %s not found", id)
	}

	return nil
}

// scanItem scans a row into a GenericItem
func (r *GenericRepository) scanItem(row *sql.Row) (*domain.GenericItem, error) {
	var item domain.GenericItem

	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Status,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("error scanning item: %w", err)
	}

	return &item, nil
}

// scanItems scans multiple rows into GenericItems
func (r *GenericRepository) scanItems(rows *sql.Rows) ([]*domain.GenericItem, error) {
	var items []*domain.GenericItem

	for rows.Next() {
		var item domain.GenericItem

		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Status,
			&item.CreatedAt,
			&item.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning item: %w", err)
		}

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return items, nil
}
