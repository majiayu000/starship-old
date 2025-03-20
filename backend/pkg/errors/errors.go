package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application-specific error
type AppError struct {
	Code    int    `json:"-"`       // HTTP status code
	Message string `json:"message"` // User-friendly error message
	Details string `json:"details"` // Detailed error information
	Err     error  `json:"-"`       // Original error
}

// Error returns the error message
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the original error
func (e *AppError) Unwrap() error {
	return e.Err
}

// StatusCode returns the HTTP status code
func (e *AppError) StatusCode() int {
	return e.Code
}

// NewBadRequest creates a new bad request error
func NewBadRequest(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Err:     err,
	}
}

// NewUnauthorized creates a new unauthorized error
func NewUnauthorized(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
		Err:     err,
	}
}

// NewForbidden creates a new forbidden error
func NewForbidden(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: message,
		Err:     err,
	}
}

// NewNotFound creates a new not found error
func NewNotFound(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
		Err:     err,
	}
}

// NewConflict creates a new conflict error
func NewConflict(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: message,
		Err:     err,
	}
}

// NewInternal creates a new internal server error
func NewInternal(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}

// WithDetails adds detail information to an AppError
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}
