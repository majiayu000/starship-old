package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/majiayu000/cc-starship/internal/core/domain"
	"github.com/majiayu000/cc-starship/internal/infrastructure/database"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

// UserRepository implements the user repository using PostgreSQL
type UserRepository struct {
	db     *database.PostgresDB
	logger *logger.Logger
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *database.PostgresDB, logger *logger.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

// FindAll finds all users
func (r *UserRepository) FindAll(ctx context.Context) ([]*domain.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password, first_name, last_name, role, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Role,
		user.Active,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// bootstrapAdminAdvisoryLock is a fixed key for serializing first-admin inserts.
const bootstrapAdminAdvisoryLock int64 = 0x73746172 // "star"

// CreateAdminIfAbsent inserts the user only when no admin exists.
// Uses a transaction advisory lock so concurrent startups cannot both succeed.
func (r *UserRepository) CreateAdminIfAbsent(ctx context.Context, user *domain.User) (bool, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("failed to begin bootstrap admin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock($1)`, bootstrapAdminAdvisoryLock); err != nil {
		return false, fmt.Errorf("failed to acquire bootstrap admin lock: %w", err)
	}

	var adminExists bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE role = 'admin')`).Scan(&adminExists); err != nil {
		return false, fmt.Errorf("failed to check for existing admin: %w", err)
	}
	if adminExists {
		if err := tx.Commit(); err != nil {
			return false, fmt.Errorf("failed to commit bootstrap admin transaction: %w", err)
		}
		return false, nil
	}

	query := `
		INSERT INTO users (id, email, password, first_name, last_name, role, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	if _, err := tx.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Role,
		user.Active,
		user.CreatedAt,
		user.UpdatedAt,
	); err != nil {
		return false, fmt.Errorf("failed to create bootstrap admin: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("failed to commit bootstrap admin transaction: %w", err)
	}
	return true, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, first_name = $2, last_name = $3, role = $4, active = $5, updated_at = $6
		WHERE id = $7
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Role,
		user.Active,
		now,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	// Update the updated_at field
	user.UpdatedAt = now

	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
