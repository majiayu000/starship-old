package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/majiayu000/cc-starship/internal/core/domain"
	"github.com/majiayu000/cc-starship/pkg/config"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

// stubAuthService satisfies ports.AuthService for route registration tests.
type stubAuthService struct{}

func (s *stubAuthService) Register(ctx context.Context, email, password, firstName, lastName string) (*domain.User, error) {
	return nil, nil
}

func (s *stubAuthService) Login(ctx context.Context, email, password string) (string, error) {
	return "", nil
}

func (s *stubAuthService) ValidateToken(ctx context.Context, token string) (*domain.User, error) {
	return &domain.User{ID: "test-user", Role: "admin"}, nil
}

// stubUserService satisfies ports.UserService so auth middleware is wired.
type stubUserService struct{}

func (s *stubUserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return nil, nil
}

func (s *stubUserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}

func (s *stubUserService) GetUsers(ctx context.Context) ([]*domain.User, error) {
	return nil, nil
}

func (s *stubUserService) CreateUser(ctx context.Context, user *domain.User) error {
	return nil
}

func (s *stubUserService) UpdateUser(ctx context.Context, user *domain.User) error {
	return nil
}

func (s *stubUserService) DeleteUser(ctx context.Context, id string) error {
	return nil
}

// stubReviewService satisfies ports.ReviewService; handlers must not be reached without auth.
type stubReviewService struct {
	reviewCalled bool
}

func (s *stubReviewService) GetByOriginalID(ctx context.Context, dataSource, originalID string) (interface{}, error) {
	return nil, nil
}

func (s *stubReviewService) GetAll(ctx context.Context, params domain.QueryParams) (*domain.PaginatedResult, error) {
	return nil, nil
}

func (s *stubReviewService) ReviewItem(ctx context.Context, dataSource, originalID string, update domain.ReviewStatusUpdate) error {
	s.reviewCalled = true
	return nil
}

func (s *stubReviewService) GetFilterOptions(ctx context.Context, dataSource string) (map[string][]string, error) {
	return nil, nil
}

func (s *stubReviewService) GetDataSources(ctx context.Context) ([]string, error) {
	return nil, nil
}

func TestUnauthenticatedReviewWriteReturns401(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	cfg := &config.Config{}
	log := logger.New(config.LoggerConfig{Level: "error", Format: "text"})
	reviewSvc := &stubReviewService{}

	RegisterRoutes(
		router,
		&stubUserService{},
		&stubAuthService{},
		nil,
		reviewSvc,
		cfg,
		log,
	)

	body := `{"status":"approved"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/review/books/items/item-1/review", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d; body=%s", http.StatusUnauthorized, rec.Code, rec.Body.String())
	}
	if reviewSvc.reviewCalled {
		t.Fatal("expected ReviewItem handler/service not to be invoked without auth")
	}
}

func TestUnauthenticatedReviewWriteReturns401WithoutUserService(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	cfg := &config.Config{}
	log := logger.New(config.LoggerConfig{Level: "error", Format: "text"})
	reviewSvc := &stubReviewService{}

	// Auth is available but user management is not — middleware must still be wired.
	RegisterRoutes(
		router,
		nil,
		&stubAuthService{},
		nil,
		reviewSvc,
		cfg,
		log,
	)

	body := `{"status":"approved"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/review/books/items/item-1/review", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d; body=%s", http.StatusUnauthorized, rec.Code, rec.Body.String())
	}
	if reviewSvc.reviewCalled {
		t.Fatal("expected ReviewItem handler/service not to be invoked without auth")
	}
}
