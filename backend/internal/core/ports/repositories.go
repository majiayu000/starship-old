package ports

import (
	"context"

	"github.com/majiayu000/cc-starship/internal/core/domain"
)

// UserRepository defines the contract for user data operations
type UserRepository interface {
	// Find users
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindAll(ctx context.Context) ([]*domain.User, error)

	// Create/Update/Delete
	Create(ctx context.Context, user *domain.User) error
	// CreateAdminIfAbsent inserts user only when no admin row exists (atomic).
	// Returns created=false when an admin already exists.
	CreateAdminIfAbsent(ctx context.Context, user *domain.User) (created bool, err error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}
