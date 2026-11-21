package errors_test

import (
	stderrors "errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/auchon/auchon-api/lib/errors"
)

type notFoundErr struct{ errors.Base }

func (notFoundErr) IsNotFound() bool { return true }

type validationErr struct{ errors.Base }

func (validationErr) IsValidationError() bool { return true }

type noAccessErr struct{ errors.Base }

func (noAccessErr) IsNoAccess() bool { return true }

type unauthorizedErr struct{ errors.Base }

func (unauthorizedErr) IsUnauthorized() bool { return true }

type duplicateErr struct{ errors.Base }

func (duplicateErr) IsDuplicate() bool { return true }

func TestCategories(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		status   int
		category string
	}{
		{"not found", notFoundErr{errors.Base{Message: "x"}}, http.StatusNotFound, errors.CategoryNotFound},
		{"validation", validationErr{errors.Base{Message: "x"}}, http.StatusBadRequest, errors.CategoryValidation},
		{"no access", noAccessErr{errors.Base{Message: "x"}}, http.StatusForbidden, errors.CategoryNoAccess},
		{"unauthorized", unauthorizedErr{errors.Base{Message: "x"}}, http.StatusUnauthorized, errors.CategoryUnauthorized},
		{"duplicate", duplicateErr{errors.Base{Message: "x"}}, http.StatusConflict, errors.CategoryDuplicate},
		{"uncategorised", stderrors.New("boom"), http.StatusInternalServerError, errors.CategoryInternal},
	}

	for _, c := range cases {
		t.Run("Given a "+c.name+" error", func(t *testing.T) {
			t.Run("When it is mapped", func(t *testing.T) {
				t.Run("Then it gets its status and category", func(t *testing.T) {
					if got := errors.HTTPStatus(c.err); got != c.status {
						t.Errorf("status = %d, want %d", got, c.status)
					}
					if got := errors.CategoryOf(c.err); got != c.category {
						t.Errorf("category = %q, want %q", got, c.category)
					}
				})
			})

			t.Run("And it is wrapped", func(t *testing.T) {
				wrapped := fmt.Errorf("context: %w", c.err)
				t.Run("Then the wrapped category is found", func(t *testing.T) {
					if got := errors.HTTPStatus(wrapped); got != c.status {
						t.Errorf("status = %d, want %d", got, c.status)
					}
				})
			})
		})
	}

	t.Run("Given no error, Then the status is 200", func(t *testing.T) {
		if got := errors.HTTPStatus(nil); got != http.StatusOK {
			t.Errorf("status = %d, want 200", got)
		}
	})

	t.Run("Given a Base, Then Error returns its message", func(t *testing.T) {
		if got := (errors.Base{Message: "box not found"}).Error(); got != "box not found" {
			t.Errorf("Error() = %q", got)
		}
	})
}
