package services

import (
	"context"
	"testing"
	"time"

	"github.com/majiayu000/cc-starship/internal/core/domain"
	"github.com/majiayu000/cc-starship/internal/infrastructure/auth"
	"github.com/majiayu000/cc-starship/pkg/config"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

type memoryUserRepo struct {
	users []*domain.User
}

func (r *memoryUserRepo) FindByID(_ context.Context, id string) (*domain.User, error) {
	for _, u := range r.users {
		if u.ID == id {
			copy := *u
			return &copy, nil
		}
	}
	return nil, errUserNotFound
}

func (r *memoryUserRepo) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			copy := *u
			return &copy, nil
		}
	}
	return nil, errUserNotFound
}

func (r *memoryUserRepo) FindAll(context.Context) ([]*domain.User, error) {
	out := make([]*domain.User, 0, len(r.users))
	for _, u := range r.users {
		copy := *u
		out = append(out, &copy)
	}
	return out, nil
}

func (r *memoryUserRepo) Create(_ context.Context, user *domain.User) error {
	copy := *user
	r.users = append(r.users, &copy)
	return nil
}

func (r *memoryUserRepo) Update(context.Context, *domain.User) error { return nil }
func (r *memoryUserRepo) Delete(context.Context, string) error         { return nil }

type userNotFoundError struct{}

func (userNotFoundError) Error() string { return "user not found" }

var errUserNotFound = userNotFoundError{}

func newTestAuthService(repo *memoryUserRepo) *AuthService {
	cfg := &config.Config{}
	cfg.Auth.JWT.SecretKey = "test-secret"
	cfg.Auth.JWT.Issuer = "test"
	cfg.Auth.JWT.ExpiryDuration = time.Hour
	log := logger.New(config.LoggerConfig{Level: "error", Format: "text"})
	return NewAuthService(repo, auth.NewJWTService(cfg), log)
}

func TestRegisterBootstrapsFirstAdmin(t *testing.T) {
	repo := &memoryUserRepo{}
	svc := newTestAuthService(repo)

	user, err := svc.Register(context.Background(), "admin@example.com", "password", "Ada", "Admin")
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if user.Role != "admin" {
		t.Fatalf("expected first registrant role admin, got %q", user.Role)
	}
}

func TestRegisterDoesNotPromoteWhenAdminExists(t *testing.T) {
	now := time.Now()
	repo := &memoryUserRepo{
		users: []*domain.User{{
			ID:        "existing-admin",
			Email:     "root@example.com",
			Role:      "admin",
			Active:    true,
			CreatedAt: now,
			UpdatedAt: now,
		}},
	}
	svc := newTestAuthService(repo)

	user, err := svc.Register(context.Background(), "user@example.com", "password", "Normal", "User")
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if user.Role != "user" {
		t.Fatalf("expected subsequent registrant role user, got %q", user.Role)
	}
}
