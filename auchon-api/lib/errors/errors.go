// Package errors is the kernel's base domain error, its categories and the error → HTTP status
// mapper. Domain errors embed Base and implement exactly one category marker.
package errors

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
	return ""
}

// HTTPStatus maps err's category to an HTTP status; 500 when it has none, 200 for nil.
func HTTPStatus(err error) int {
	return 0
}
