package domain

import (
	"fmt"

	"github.com/auchon/auchon-api/lib/errors"
)

type LocationCannotContainItselfError struct{ errors.Base }

func NewLocationCannotContainItselfError(id string) LocationCannotContainItselfError {
	return LocationCannotContainItselfError{errors.Base{Message: fmt.Sprintf("location %s cannot go inside itself or a location it contains", id), Context: map[string]any{"locationId": id}}}
}

func (LocationCannotContainItselfError) IsValidationError() bool { return true }
