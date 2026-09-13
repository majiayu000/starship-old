package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/majiayu000/cc-starship/internal/core/domain"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

type mockUserService struct {
	users   map[string]*domain.User
	updated *domain.User
}

func (m *mockUserService) listUsers() []*domain.User {
	out := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		copy := *user
		out = append(out, &copy)
	}
	return out
}

func (m *mockUserService) GetUserByID(_ context.Context, id string) (*domain.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, errNotFound
	}
	copy := *user
	return &copy, nil
}

func (m *mockUserService) GetUserByEmail(context.Context, string) (*domain.User, error) {
	return nil, nil
}

func (m *mockUserService) GetUsers(context.Context) ([]*domain.User, error) {
	return m.listUsers(), nil
}

func (m *mockUserService) CreateUser(context.Context, *domain.User) error {
	return nil
}

func (m *mockUserService) UpdateUser(_ context.Context, user *domain.User) error {
	m.updated = user
	m.users[user.ID] = user
	return nil
}

func (m *mockUserService) DeleteUser(context.Context, string) error {
	return nil
}

type notFoundError struct{}

func (notFoundError) Error() string { return "not found" }

var errNotFound = notFoundError{}

func newTestUser(id, email, role string, active bool) *domain.User {
	now := time.Now()
	return &domain.User{
		ID:        id,
		Email:     email,
		FirstName: "First",
		LastName:  "Last",
		Role:      role,
		Active:    active,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func setupUpdateRouter(handler *UserHandler, caller *domain.User) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.PUT("/api/v1/users/:id", func(c *gin.Context) {
		if caller != nil {
			c.Set("user", caller)
		}
		handler.UpdateUser(c)
	})
	return r
}

func TestUpdateUser_NonAdminSelfProfileUpdateSucceedsWithoutRoleChange(t *testing.T) {
	target := newTestUser("user-1", "user@example.com", "user", true)
	svc := &mockUserService{users: map[string]*domain.User{"user-1": target}}
	handler := NewUserHandler(svc, &logger.Logger{})
	caller := newTestUser("user-1", "user@example.com", "user", true)
	router := setupUpdateRouter(handler, caller)

	body := `{"email":"user@example.com","firstName":"Updated","lastName":"Name","role":"admin","active":false}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/user-1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", w.Code, w.Body.String())
	}
	if svc.updated == nil {
		t.Fatal("expected user to be updated")
	}
	if svc.updated.FirstName != "Updated" || svc.updated.LastName != "Name" {
		t.Fatalf("expected profile fields updated, got firstName=%q lastName=%q", svc.updated.FirstName, svc.updated.LastName)
	}
	if svc.updated.Role != "user" {
		t.Fatalf("expected role to remain user, got %q", svc.updated.Role)
	}
	if !svc.updated.Active {
		t.Fatal("expected active to remain true for non-admin update")
	}
}

func TestUpdateUser_NonAdminUpdateOtherUserReturns403(t *testing.T) {
	target := newTestUser("user-2", "other@example.com", "user", true)
	svc := &mockUserService{users: map[string]*domain.User{"user-2": target}}
	handler := NewUserHandler(svc, &logger.Logger{})
	caller := newTestUser("user-1", "user@example.com", "user", true)
	router := setupUpdateRouter(handler, caller)

	body := `{"email":"other@example.com","firstName":"Hacked","lastName":"User"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/user-2", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d body=%s", w.Code, w.Body.String())
	}
	if svc.updated != nil {
		t.Fatal("expected no update for forbidden request")
	}
}

func TestUpdateUser_NonAdminRoleEscalationLeavesRoleUnchanged(t *testing.T) {
	target := newTestUser("user-1", "user@example.com", "user", true)
	svc := &mockUserService{users: map[string]*domain.User{"user-1": target}}
	handler := NewUserHandler(svc, &logger.Logger{})
	caller := newTestUser("user-1", "user@example.com", "user", true)
	router := setupUpdateRouter(handler, caller)

	body := `{"email":"user@example.com","firstName":"First","lastName":"Last","role":"admin"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/user-1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	data, _ := resp["data"].(map[string]interface{})
	if data["role"] != "user" {
		t.Fatalf("expected response role user, got %v", data["role"])
	}
	if svc.updated.Role != "user" {
		t.Fatalf("expected persisted role user, got %q", svc.updated.Role)
	}
}

func TestUpdateUser_AdminCanUpdateOtherUserIncludingRole(t *testing.T) {
	target := newTestUser("user-2", "other@example.com", "user", true)
	svc := &mockUserService{users: map[string]*domain.User{"user-2": target}}
	handler := NewUserHandler(svc, &logger.Logger{})
	caller := newTestUser("admin-1", "admin@example.com", "admin", true)
	router := setupUpdateRouter(handler, caller)

	body := `{"email":"other@example.com","firstName":"Promoted","lastName":"User","role":"admin","active":false}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/user-2", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", w.Code, w.Body.String())
	}
	if svc.updated == nil {
		t.Fatal("expected user to be updated")
	}
	if svc.updated.Role != "admin" {
		t.Fatalf("expected admin to set role, got %q", svc.updated.Role)
	}
	if svc.updated.Active {
		t.Fatal("expected admin to set active=false")
	}
	if svc.updated.FirstName != "Promoted" {
		t.Fatalf("expected firstName Promoted, got %q", svc.updated.FirstName)
	}
}

func TestUpdateUser_RejectsSoleAdminSelfDemotion(t *testing.T) {
	admin := newTestUser("admin-1", "admin@example.com", "admin", true)
	svc := &mockUserService{users: map[string]*domain.User{"admin-1": admin}}
	handler := NewUserHandler(svc, &logger.Logger{})
	router := setupUpdateRouter(handler, admin)

	body := `{"email":"admin@example.com","firstName":"First","lastName":"Last","role":"user"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/admin-1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d body=%s", w.Code, w.Body.String())
	}
	if svc.updated != nil {
		t.Fatal("expected sole admin demotion to be rejected")
	}
}

func TestUpdateUser_AllowsAdminDemotionWhenAnotherAdminExists(t *testing.T) {
	admin1 := newTestUser("admin-1", "admin1@example.com", "admin", true)
	admin2 := newTestUser("admin-2", "admin2@example.com", "admin", true)
	svc := &mockUserService{users: map[string]*domain.User{
		"admin-1": admin1,
		"admin-2": admin2,
	}}
	handler := NewUserHandler(svc, &logger.Logger{})
	router := setupUpdateRouter(handler, admin1)

	body := `{"email":"admin1@example.com","firstName":"First","lastName":"Last","role":"user"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/admin-1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", w.Code, w.Body.String())
	}
	if svc.updated == nil || svc.updated.Role != "user" {
		t.Fatalf("expected demotion when another admin exists, updated=%v", svc.updated)
	}
}
