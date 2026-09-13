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

// CreateAdminIfAbsent inserts or reactivates an admin only when no active admin
// exists. Uses a transaction advisory lock so concurrent startups cannot both succeed.
// Inactive-only admin rows do not block recovery: matching bootstrap email is
// reactivated/promoted; otherwise a new admin is inserted.
func (r *UserRepository) CreateAdminIfAbsent(ctx context.Context, user *domain.User) (bool, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("failed to begin bootstrap admin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock($1)`, bootstrapAdminAdvisoryLock); err != nil {
		return false, fmt.Errorf("failed to acquire bootstrap admin lock: %w", err)
	}

	var activeAdminExists bool
	if err := tx.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM users WHERE role = 'admin' AND active = true)
	`).Scan(&activeAdminExists); err != nil {
		return false, fmt.Errorf("failed to check for existing active admin: %w", err)
	}
	if activeAdminExists {
		if err := tx.Commit(); err != nil {
			return false, fmt.Errorf("failed to commit bootstrap admin transaction: %w", err)
		}
		return false, nil
	}

	var existingID string
	err = tx.QueryRowContext(ctx, `
		SELECT id FROM users WHERE email = $1 FOR UPDATE
	`, user.Email).Scan(&existingID)
	if err == nil {
		now := time.Now()
		if _, err := tx.ExecContext(ctx, `
			UPDATE users
			SET password = $1, first_name = $2, last_name = $3, role = 'admin', active = true, updated_at = $4
			WHERE id = $5
		`, user.Password, user.FirstName, user.LastName, now, existingID); err != nil {
			return false, fmt.Errorf("failed to reactivate bootstrap admin: %w", err)
		}
		user.ID = existingID
		user.Role = "admin"
		user.Active = true
		user.UpdatedAt = now
		if err := tx.Commit(); err != nil {
			return false, fmt.Errorf("failed to commit bootstrap admin transaction: %w", err)
		}
		return true, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf("failed to look up bootstrap admin email: %w", err)
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

// Update updates an existing user. Losing the last active admin (demotion or
// deactivation) is rejected atomically. Admin locks are always acquired in
// ORDER BY id before mutating, avoiding deadlocks between concurrent demotions.
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin user update transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var currentRole string
	var currentActive bool
	err = tx.QueryRowContext(ctx, `
		SELECT role, active FROM users WHERE id = $1
	`, user.ID).Scan(&currentRole, &currentActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to load user for update: %w", err)
	}

	// Serialize admin capability loss with consistent lock order (ORDER BY id)
	// before any exclusive target lock, so concurrent demotions cannot deadlock.
	needsAdminLockSet := currentRole == "admin" && (user.Role != "admin" || !user.Active)

	if needsAdminLockSet {
		rows, err := tx.QueryContext(ctx, `
			SELECT id, active FROM users WHERE role = 'admin' ORDER BY id FOR UPDATE
		`)
		if err != nil {
			return fmt.Errorf("failed to lock admin users: %w", err)
		}
		activeAdminCount := 0
		targetSeen := false
		targetActive := false
		for rows.Next() {
			var id string
			var active bool
			if err := rows.Scan(&id, &active); err != nil {
				rows.Close()
				return fmt.Errorf("failed to scan admin row: %w", err)
			}
			if active {
				activeAdminCount++
			}
			if id == user.ID {
				targetSeen = true
				targetActive = active
			}
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return fmt.Errorf("error iterating admin users: %w", err)
		}
		rows.Close()

		losingUsableAdmin := targetSeen && targetActive && (user.Role != "admin" || !user.Active)
		if losingUsableAdmin && activeAdminCount <= 1 {
			return domain.ErrCannotDemoteLastAdmin
		}

		if !targetSeen {
			err = tx.QueryRowContext(ctx, `
				SELECT role, active FROM users WHERE id = $1 FOR UPDATE
			`, user.ID).Scan(&currentRole, &currentActive)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return fmt.Errorf("user not found")
				}
				return fmt.Errorf("failed to lock user for update: %w", err)
			}
		}
	} else {
		err = tx.QueryRowContext(ctx, `
			SELECT role, active FROM users WHERE id = $1 FOR UPDATE
		`, user.ID).Scan(&currentRole, &currentActive)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("user not found")
			}
			return fmt.Errorf("failed to lock user for update: %w", err)
		}
	}

	now := time.Now()
	result, err := tx.ExecContext(ctx, `
		UPDATE users
		SET email = $1, first_name = $2, last_name = $3, role = $4, active = $5, updated_at = $6
		WHERE id = $7
	`,
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

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit user update: %w", err)
	}

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
