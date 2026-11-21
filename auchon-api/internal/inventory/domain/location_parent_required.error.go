package domain

import (
	"fmt"

	"github.com/auchon/auchon-api/lib/errors"
)

type LocationParentRequiredError struct{ errors.Base }

func NewLocationParentRequiredError(id string) LocationParentRequiredError {
	return LocationParentRequiredError{errors.Base{Message: fmt.Sprintf("location %s must be inside another location", id), Context: map[string]any{"locationId": id}}}
}

func (LocationParentRequiredError) IsValidationError() bool { return true }
