package domain

import (
	"time"
)

// Status constants for generic items
const (
	StatusDraft    = "draft"
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusArchived = "archived"
)

// BaseEntity represents common fields for any entity
type BaseEntity struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// GenericItem represents a generic item with common fields
type GenericItem struct {
	BaseEntity
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// StatusUpdate represents a request to update an item's status
type StatusUpdate struct {
	Status    string `json:"status"`
	Comment   string `json:"comment,omitempty"`
	UpdatedBy string `json:"updatedBy,omitempty"`
}

// PaginatedResult represents a paginated result set
type PaginatedResult struct {
	Items      interface{} `json:"items"`
	TotalItems int64       `json:"totalItems"`
	TotalPages int         `json:"totalPages"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
}

// QueryParams represents query parameters for filtering and pagination
type QueryParams struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	Status     string `json:"status,omitempty"`     // Filter by status
	SearchTerm string `json:"searchTerm,omitempty"` // Search term for content
	SortBy     string `json:"sortBy,omitempty"`     // Field to sort by
	SortOrder  string `json:"sortOrder,omitempty"`  // 'asc' or 'desc'
}

// IsValidStatus checks if a status is valid
func IsValidStatus(status string) bool {
	return status == StatusDraft ||
		status == StatusActive ||
		status == StatusInactive ||
		status == StatusArchived
}
