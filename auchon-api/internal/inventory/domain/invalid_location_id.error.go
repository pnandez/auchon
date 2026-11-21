package domain

import (
	"fmt"

	"github.com/auchon/auchon-api/lib/errors"
)

type InvalidLocationIDError struct{ errors.Base }

func NewInvalidLocationIDError(id string) InvalidLocationIDError {
	return InvalidLocationIDError{errors.Base{Message: fmt.Sprintf("%q is not a valid location id", id), Context: map[string]any{"locationId": id}}}
}

func (InvalidLocationIDError) IsValidationError() bool { return true }
