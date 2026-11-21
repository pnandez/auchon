// Package errors is the kernel's base domain error, its categories and the error → HTTP status
// mapper. Domain errors embed Base and implement exactly one category marker.
package errors

import (
	stderrors "errors"
	"net/http"
)

// Base is embedded in every domain error.
type Base struct {
	Message string
	Context map[string]any
}

func (b Base) Error() string { return b.Message }

// Category markers. A domain error implements exactly one.
type (
	NotFound        interface{ IsNotFound() bool }
	ValidationError interface{ IsValidationError() bool }
	NoAccess        interface{ IsNoAccess() bool }
	Unauthorized    interface{ IsUnauthorized() bool }
	Duplicate       interface{ IsDuplicate() bool }
)

// Category names, used by the sync push result.
const (
	CategoryNotFound     = "not_found"
	CategoryValidation   = "validation"
	CategoryNoAccess     = "no_access"
	CategoryUnauthorized = "unauthorized"
	CategoryDuplicate    = "duplicate"
	CategoryInternal     = "internal"
)

// CategoryOf returns the category of err (searching wrapped errors), or CategoryInternal.
func CategoryOf(err error) string {
	var (
		notFound     NotFound
		validation   ValidationError
		noAccess     NoAccess
		unauthorized Unauthorized
		duplicate    Duplicate
	)
	switch {
	case stderrors.As(err, &notFound):
		return CategoryNotFound
	case stderrors.As(err, &validation):
		return CategoryValidation
	case stderrors.As(err, &noAccess):
		return CategoryNoAccess
	case stderrors.As(err, &unauthorized):
		return CategoryUnauthorized
	case stderrors.As(err, &duplicate):
		return CategoryDuplicate
	default:
		return CategoryInternal
	}
}

var statusByCategory = map[string]int{
	CategoryNotFound:     http.StatusNotFound,
	CategoryValidation:   http.StatusBadRequest,
	CategoryNoAccess:     http.StatusForbidden,
	CategoryUnauthorized: http.StatusUnauthorized,
	CategoryDuplicate:    http.StatusConflict,
	CategoryInternal:     http.StatusInternalServerError,
}

// HTTPStatus maps err's category to an HTTP status; 500 when it has none, 200 for nil.
func HTTPStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}
	return statusByCategory[CategoryOf(err)]
}
