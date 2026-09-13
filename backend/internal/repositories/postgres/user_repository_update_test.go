package postgres

import (
	"testing"

	"github.com/majiayu000/cc-starship/internal/core/domain"
)

func TestResolvePrivilegedUpdateFields_PreservesLockedWhenNotProvided(t *testing.T) {
	user := &domain.User{
		Role:   "user",  // stale pre-lock value
		Active: true,    // stale pre-lock value
	}
	role, active := resolvePrivilegedUpdateFields("admin", false, user)
	if role != "admin" || active {
		t.Fatalf("expected locked admin/false, got role=%q active=%v", role, active)
	}
}

func TestResolvePrivilegedUpdateFields_AppliesExplicitRoleAndActive(t *testing.T) {
	user := &domain.User{
		Role:           "user",
		Active:         false,
		RoleProvided:   true,
		ActiveProvided: true,
	}
	role, active := resolvePrivilegedUpdateFields("admin", true, user)
	if role != "user" || active {
		t.Fatalf("expected explicit user/false, got role=%q active=%v", role, active)
	}
}

func TestResolvePrivilegedUpdateFields_AppliesOnlyProvidedFields(t *testing.T) {
	user := &domain.User{
		Role:         "admin",
		Active:       true, // stale; must not overwrite locked false
		RoleProvided: true,
	}
	role, active := resolvePrivilegedUpdateFields("user", false, user)
	if role != "admin" {
		t.Fatalf("expected explicit role admin, got %q", role)
	}
	if active {
		t.Fatal("expected locked active=false preserved when ActiveProvided=false")
	}
}
