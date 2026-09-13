package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/majiayu000/cc-starship/internal/core/domain"
	"github.com/majiayu000/cc-starship/internal/infrastructure/database"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

// ReviewRepository implements the review repository using PostgreSQL
type ReviewRepository struct {
	db     *database.PostgresDB
	logger *logger.Logger
	// Map of data sources to their corresponding tables
	sourceToTable map[string]string
}

// NewReviewRepository creates a new ReviewRepository
func NewReviewRepository(db *database.PostgresDB, logger *logger.Logger) *ReviewRepository {
	return &ReviewRepository{
		db:     db,
		logger: logger,
		sourceToTable: map[string]string{
			domain.DataSourceSATOneprep: "t_math_sat_temp_classify_oneprep",
			domain.DataSourceSATIXL:     "t_math_sat_temp_classify_ixl", // Update with the actual table name
		},
	}
}

// FindByOriginalID retrieves a reviewable item by its original ID and data source
func (r *ReviewRepository) FindByOriginalID(ctx context.Context, dataSource, originalID string) (interface{}, error) {
	tableName, err := r.getTableNameForDataSource(dataSource)
	if err != nil {
		return nil, err
	}

	// Construct the query based on the data source
	query := fmt.Sprintf(`SELECT * FROM %s WHERE original_id = $1`, tableName)

	row := r.db.QueryRowContext(ctx, query, originalID)

	switch dataSource {
	case domain.DataSourceSATOneprep:
		return r.scanSATOneprepItem(row)
	case domain.DataSourceSATIXL:
		return r.scanSATIXLItem(row)
	default:
		return nil, fmt.Errorf("unsupported data source: %s", dataSource)
	}
}

// FindAll retrieves all reviewable items with pagination and filtering
func (r *ReviewRepository) FindAll(ctx context.Context, params domain.QueryParams) (*domain.PaginatedResult, error) {
	tableName, err := r.getTableNameForDataSource(params.DataSource)
	if err != nil {
		return nil, err
	}

	// Construct the base query
	baseQuery := fmt.Sprintf("FROM %s", tableName)

	// Add WHERE clauses for filtering
	whereClause, args := r.buildWhereClause(params)
	if whereClause != "" {
		baseQuery += " WHERE " + whereClause
	}

	// Count total items for pagination
	countQuery := fmt.Sprintf("SELECT COUNT(*) %s", baseQuery)
	var totalItems int64
	err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("error counting total items: %w", err)
	}

	// Resolve ORDER BY from an allowlist only — never interpolate raw SortBy.
	orderBy, err := domain.ResolveSortColumn(params.DataSource, params.SortBy)
	if err != nil {
		return nil, err
	}

	sortOrder := "ASC"
	if strings.ToLower(params.SortOrder) == "desc" {
		sortOrder = "DESC"
	}

	// Ensure pagination parameters are valid
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	offset := (params.Page - 1) * params.PageSize

	dataQuery := fmt.Sprintf("SELECT * %s ORDER BY %s %s LIMIT %d OFFSET %d",
		baseQuery, orderBy, sortOrder, params.PageSize, offset)

	// Execute the query
	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying data: %w", err)
	}
	defer rows.Close()

	// Scan the results
	var items interface{}
	switch params.DataSource {
	case domain.DataSourceSATOneprep:
		items, err = r.scanMultipleSATOneprepItems(rows)
	case domain.DataSourceSATIXL:
		items, err = r.scanMultipleSATIXLItems(rows)
	default:
		return nil, fmt.Errorf("unsupported data source: %s", params.DataSource)
	}

	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalItems) / float64(params.PageSize)))

	return &domain.PaginatedResult{
		Items:      items,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Page:       params.Page,
		PageSize:   params.PageSize,
		DataSource: params.DataSource,
	}, nil
}

// UpdateReviewStatus updates the review status of a reviewable item
func (r *ReviewRepository) UpdateReviewStatus(ctx context.Context, dataSource, originalID string, update domain.ReviewStatusUpdate) error {
	tableName, err := r.getTableNameForDataSource(dataSource)
	if err != nil {
		return err
	}

	// Validate the review status
	if !domain.IsValidReviewStatus(update.Status) {
		return fmt.Errorf("invalid review status: %s", update.Status)
	}

	// Update only the review_status field
	query := fmt.Sprintf(`
		UPDATE %s
		SET review_status = $1
		WHERE original_id = $2
	`, tableName)

	result, err := r.db.ExecContext(ctx, query, update.Status, originalID)
	if err != nil {
		return fmt.Errorf("error updating review status: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no item found with original_id %s", originalID)
	}

	return nil
}

// GetFilterOptions retrieves available options for filtering
func (r *ReviewRepository) GetFilterOptions(ctx context.Context, dataSource string) (map[string][]string, error) {
	tableName, err := r.getTableNameForDataSource(dataSource)
	if err != nil {
		return nil, err
	}

	options := make(map[string][]string)

	// Get review statuses
	options["reviewStatus"] = []string{
		domain.ReviewStatusPending,
		domain.ReviewStatusApproved,
		domain.ReviewStatusRejected,
	}

	// Get domains (if applicable)
	if dataSource == domain.DataSourceSATOneprep {
		domains, err := r.getDistinctValues(ctx, tableName, "domain")
		if err != nil {
			return nil, err
		}
		options["domain"] = domains

		// Get skills (if applicable)
		skills, err := r.getDistinctValues(ctx, tableName, "skill")
		if err != nil {
			return nil, err
		}
		options["skill"] = skills
	}

	return options, nil
}

// GetDataSources returns a list of available data sources
func (r *ReviewRepository) GetDataSources(ctx context.Context) ([]string, error) {
	sources := []string{
		domain.DataSourceSATOneprep,
		domain.DataSourceSATIXL,
	}
	return sources, nil
}

// Helper methods

// getTableNameForDataSource returns the table name for a given data source
func (r *ReviewRepository) getTableNameForDataSource(dataSource string) (string, error) {
	tableName, ok := r.sourceToTable[dataSource]
	if !ok {
		return "", fmt.Errorf("unsupported data source: %s", dataSource)
	}
	return tableName, nil
}

// buildWhereClause constructs the WHERE clause for the query based on the params
func (r *ReviewRepository) buildWhereClause(params domain.QueryParams) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	paramIdx := 1

	// Add review status filter
	if params.ReviewStatus != "" {
		conditions = append(conditions, fmt.Sprintf("review_status = $%d", paramIdx))
		args = append(args, params.ReviewStatus)
		paramIdx++
	}

	// Add domain filter (if applicable)
	if params.Domain != "" {
		conditions = append(conditions, fmt.Sprintf("domain = $%d", paramIdx))
		args = append(args, params.Domain)
		paramIdx++
	}

	// Add skill filter (if applicable)
	if params.Skill != "" {
		conditions = append(conditions, fmt.Sprintf("skill = $%d", paramIdx))
		args = append(args, params.Skill)
		paramIdx++
	}

	// Add search term filter
	if params.SearchTerm != "" {
		searchPattern := "%" + params.SearchTerm + "%"

		var searchConditions []string
		searchTables := map[string][]string{
			domain.DataSourceSATOneprep: {"question_set", "subject", "domain", "skill", "knowledge_point"},
			domain.DataSourceSATIXL:     {"knowledge_point", "skill"},
		}

		searchColumns, ok := searchTables[params.DataSource]
		if ok {
			for _, col := range searchColumns {
				searchConditions = append(searchConditions, fmt.Sprintf("%s ILIKE $%d", col, paramIdx))
			}
			args = append(args, searchPattern)
			paramIdx++

			if len(searchConditions) > 0 {
				conditions = append(conditions, "("+strings.Join(searchConditions, " OR ")+")")
			}
		}
	}

	if len(conditions) > 0 {
		return strings.Join(conditions, " AND "), args
	}
	return "", args
}

// getDistinctValues retrieves distinct values for a column
func (r *ReviewRepository) getDistinctValues(ctx context.Context, tableName, columnName string) ([]string, error) {
	query := fmt.Sprintf("SELECT DISTINCT %s FROM %s WHERE %s IS NOT NULL ORDER BY %s",
		columnName, tableName, columnName, columnName)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying distinct values: %w", err)
	}
	defer rows.Close()

	var values []string
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, fmt.Errorf("error scanning value: %w", err)
		}
		values = append(values, value)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return values, nil
}

// scanSATOneprepItem scans a row into a SATOneprepItem
func (r *ReviewRepository) scanSATOneprepItem(row *sql.Row) (*domain.SATOneprepItem, error) {
	var item domain.SATOneprepItem

	// Scan only the columns that exist in the database table (15 columns)
	err := row.Scan(
		&item.URL,
		&item.QuestionSet,
		&item.Subject,
		&item.Difficulty,
		&item.Domain,
		&item.Skill,
		&item.QuestionType,
		&item.QuestionContent,
		&item.Explanation,
		&item.Answer,
		&item.Options,
		&item.QuestionID,
		&item.KnowledgePoint,
		&item.OriginalID,
		&item.ReviewStatus,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("error scanning item: %w", err)
	}

	// Set default values for missing fields
	item.ReviewComment = ""
	item.ReviewedAt = nil
	item.ReviewedBy = ""

	return &item, nil
}

// scanSATIXLItem scans a row into a SATIXLItem
func (r *ReviewRepository) scanSATIXLItem(row *sql.Row) (*domain.SATIXLItem, error) {
	// Implement this method based on the actual table structure
	// This is just a placeholder
	return &domain.SATIXLItem{}, nil
}

// scanMultipleSATOneprepItems scans multiple rows into SATOneprepItems
func (r *ReviewRepository) scanMultipleSATOneprepItems(rows *sql.Rows) ([]*domain.SATOneprepItem, error) {
	var items []*domain.SATOneprepItem

	for rows.Next() {
		var item domain.SATOneprepItem

		// Scan only the columns that exist in the database table (15 columns)
		err := rows.Scan(
			&item.URL,
			&item.QuestionSet,
			&item.Subject,
			&item.Difficulty,
			&item.Domain,
			&item.Skill,
			&item.QuestionType,
			&item.QuestionContent,
			&item.Explanation,
			&item.Answer,
			&item.Options,
			&item.QuestionID,
			&item.KnowledgePoint,
			&item.OriginalID,
			&item.ReviewStatus,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning item: %w", err)
		}

		// Set default values for missing fields
		item.ReviewComment = ""
		item.ReviewedAt = nil
		item.ReviewedBy = ""

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return items, nil
}

// scanMultipleSATIXLItems scans multiple rows into SATIXLItems
func (r *ReviewRepository) scanMultipleSATIXLItems(rows *sql.Rows) ([]*domain.SATIXLItem, error) {
	// Implement this method based on the actual table structure
	// This is just a placeholder
	return []*domain.SATIXLItem{}, nil
}
