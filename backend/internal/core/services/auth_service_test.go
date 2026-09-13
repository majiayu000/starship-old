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

func (r *memoryUserRepo) CreateAdminIfAbsent(_ context.Context, user *domain.User) (bool, error) {
	for _, existing := range r.users {
		if existing.Role == "admin" {
			return false, nil
		}
	}
	copy := *user
	r.users = append(r.users, &copy)
	return true, nil
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

func TestRegisterNeverPromotesToAdmin(t *testing.T) {
	repo := &memoryUserRepo{}
	svc := newTestAuthService(repo)

	user, err := svc.Register(context.Background(), "first@example.com", "password", "First", "User")
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if user.Role != "user" {
		t.Fatalf("expected public registrant role user, got %q", user.Role)
	}
}

func TestEnsureBootstrapAdminCreatesFirstAdmin(t *testing.T) {
	repo := &memoryUserRepo{}
	svc := newTestAuthService(repo)

	err := svc.EnsureBootstrapAdmin(context.Background(), config.BootstrapAdminConfig{
		Email:     "admin@example.com",
		Password:  "password",
		FirstName: "Ada",
		LastName:  "Admin",
	})
	if err != nil {
		t.Fatalf("EnsureBootstrapAdmin returned error: %v", err)
	}
	if len(repo.users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(repo.users))
	}
	if repo.users[0].Role != "admin" {
		t.Fatalf("expected bootstrap role admin, got %q", repo.users[0].Role)
	}
	if repo.users[0].Email != "admin@example.com" {
		t.Fatalf("expected bootstrap email admin@example.com, got %q", repo.users[0].Email)
	}
}

func TestEnsureBootstrapAdminSkipsWhenAdminExists(t *testing.T) {
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

	err := svc.EnsureBootstrapAdmin(context.Background(), config.BootstrapAdminConfig{
		Email:    "other-admin@example.com",
		Password: "password",
	})
	if err != nil {
		t.Fatalf("EnsureBootstrapAdmin returned error: %v", err)
	}
	if len(repo.users) != 1 {
		t.Fatalf("expected no additional user, got %d", len(repo.users))
	}
}

func TestEnsureBootstrapAdminNoopWithoutCredentials(t *testing.T) {
	repo := &memoryUserRepo{}
	svc := newTestAuthService(repo)

	err := svc.EnsureBootstrapAdmin(context.Background(), config.BootstrapAdminConfig{})
	if err != nil {
		t.Fatalf("EnsureBootstrapAdmin returned error: %v", err)
	}
	if len(repo.users) != 0 {
		t.Fatalf("expected no users when bootstrap credentials empty, got %d", len(repo.users))
	}
}

func TestEnsureBootstrapAdminRejectsInvalidEmail(t *testing.T) {
	repo := &memoryUserRepo{}
	svc := newTestAuthService(repo)

	err := svc.EnsureBootstrapAdmin(context.Background(), config.BootstrapAdminConfig{
		Email:    "not-an-email",
		Password: "password",
	})
	if err == nil {
		t.Fatal("expected error for invalid bootstrap email")
	}
	if len(repo.users) != 0 {
		t.Fatalf("expected no users inserted for invalid email, got %d", len(repo.users))
	}
}
