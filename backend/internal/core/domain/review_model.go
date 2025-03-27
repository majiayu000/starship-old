package domain

import (
	"encoding/json"
	"time"
)

// ReviewStatus constants
const (
	ReviewStatusPending  = "pending"
	ReviewStatusApproved = "approved"
	ReviewStatusRejected = "rejected"
)

// DataSource constants for different data types
const (
	DataSourceSATOneprep = "sat_oneprep"
	DataSourceSATIXL     = "sat_ixl"
	// Add more data sources as needed
)

// ReviewableItem represents common fields for any item that can be reviewed
type ReviewableItem struct {
	OriginalID    string     `json:"originalId"`
	ReviewStatus  string     `json:"reviewStatus"` // 'pending', 'approved', 'rejected'
	ReviewComment string     `json:"reviewComment,omitempty"`
	ReviewedAt    *time.Time `json:"reviewedAt,omitempty"`
	ReviewedBy    string     `json:"reviewedBy,omitempty"`
}

// SATOneprepItem represents a SAT question from the oneprep source
type SATOneprepItem struct {
	ReviewableItem
	URL             string          `json:"url"`
	QuestionSet     string          `json:"questionSet"`
	Subject         string          `json:"subject"`
	Difficulty      string          `json:"difficulty"`
	Domain          string          `json:"domain"`
	Skill           string          `json:"skill"`
	QuestionType    string          `json:"questionType"`
	QuestionContent json.RawMessage `json:"questionContent"`
	Explanation     json.RawMessage `json:"explanation"`
	Answer          json.RawMessage `json:"answer"`
	Options         json.RawMessage `json:"options"`
	QuestionID      string          `json:"questionId"`
	KnowledgePoint  string          `json:"knowledgePoint"`
}

// SATIXLItem represents a SAT question from the IXL source
// Note: Update the fields based on the actual structure of the IXL data
type SATIXLItem struct {
	ReviewableItem
	// Add IXL specific fields here based on the actual data structure
	QuestionID     string          `json:"questionId"`
	Content        json.RawMessage `json:"content"`
	Answer         json.RawMessage `json:"answer"`
	Skill          string          `json:"skill"`
	KnowledgePoint string          `json:"knowledgePoint"`
	// Add more fields as needed
}

// ReviewStatusUpdate represents a request to update the review status
type ReviewStatusUpdate struct {
	Status     string `json:"status"` // 'pending', 'approved', 'rejected'
	Comment    string `json:"comment,omitempty"`
	ReviewedBy string `json:"reviewedBy,omitempty"`
}

// PaginatedResult represents a paginated result set
type PaginatedResult struct {
	Items      interface{} `json:"items"`
	TotalItems int64       `json:"totalItems"`
	TotalPages int         `json:"totalPages"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	DataSource string      `json:"dataSource"`
}

// QueryParams represents query parameters for filtering and pagination
type QueryParams struct {
	DataSource   string `json:"dataSource"` // Which data source to query
	Page         int    `json:"page"`
	PageSize     int    `json:"pageSize"`
	ReviewStatus string `json:"reviewStatus,omitempty"` // Filter by review status
	Domain       string `json:"domain,omitempty"`       // Filter by domain
	Skill        string `json:"skill,omitempty"`        // Filter by skill
	SearchTerm   string `json:"searchTerm,omitempty"`   // Search term for content
	SortBy       string `json:"sortBy,omitempty"`       // Field to sort by
	SortOrder    string `json:"sortOrder,omitempty"`    // 'asc' or 'desc'
}

// IsValidReviewStatus checks if a review status is valid
func IsValidReviewStatus(status string) bool {
	return status == ReviewStatusPending ||
		status == ReviewStatusApproved ||
		status == ReviewStatusRejected
}

// IsValidDataSource checks if a data source is valid
func IsValidDataSource(source string) bool {
	return source == DataSourceSATOneprep ||
		source == DataSourceSATIXL
	// Update as more data sources are added
}
