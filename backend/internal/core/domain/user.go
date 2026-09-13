package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrCannotDemoteLastAdmin is returned when an update would leave the system
// with no administrator.
var ErrCannotDemoteLastAdmin = errors.New("cannot demote the sole administrator")

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
