package domain

import (
	"fmt"

	"github.com/auchon/auchon-api/lib/errors"
)

type ParentLocationNotFoundError struct{ errors.Base }

func NewParentLocationNotFoundError(id string) ParentLocationNotFoundError {
	return ParentLocationNotFoundError{errors.Base{Message: fmt.Sprintf("location %s to put it in not found", id), Context: map[string]any{"locationId": id}}}
}

func (ParentLocationNotFoundError) IsNotFound() bool { return true }
