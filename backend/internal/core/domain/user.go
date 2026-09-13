package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrCannotDemoteLastAdmin is returned when an update or delete would leave the
// system with no active administrator (role demotion, deactivation, or deletion).
var ErrCannotDemoteLastAdmin = errors.New("cannot demote, deactivate, or delete the sole active administrator")

// ErrUserNotFound indicates the user row is absent (invalid session), not a DB outage.
var ErrUserNotFound = errors.New("user not found")

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never expose password in JSON
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// RoleProvided / ActiveProvided mark intentional privileged-field writes for Update.
	// When false, the repository preserves Role/Active observed under row lock so a
	// concurrent admin deactivate/role change is not overwritten by a stale pre-lock read.
	RoleProvided   bool `json:"-"`
	ActiveProvided bool `json:"-"`
}

// NewUser creates a new user with default values
func NewUser(email, password, firstName, lastName string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  password, // Note: Password should be hashed before storage
		FirstName: firstName,
		LastName:  lastName,
		Role:      "user", // Default role
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}
