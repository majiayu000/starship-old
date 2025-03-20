package ports

import (
	"context"

	"github.com/majiayu000/cc-starship/internal/core/domain"
)

// UserService defines the contract for user operations
type UserService interface {
	// User operations
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUsers(ctx context.Context) ([]*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id string) error
}

// AuthService defines the contract for authentication operations
type AuthService interface {
	// Authentication
	Register(ctx context.Context, email, password, firstName, lastName string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, error) // returns JWT token
	ValidateToken(ctx context.Context, token string) (*domain.User, error)
}
