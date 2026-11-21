package domain

import (
	"fmt"

	"github.com/auchon/auchon-api/lib/errors"
)

type LocationAlreadyExistsError struct{ errors.Base }

func NewLocationAlreadyExistsError(id string) LocationAlreadyExistsError {
	return LocationAlreadyExistsError{errors.Base{Message: fmt.Sprintf("location %s already exists", id), Context: map[string]any{"locationId": id}}}
}

func (LocationAlreadyExistsError) IsDuplicate() bool { return true }
