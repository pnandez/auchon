package domain

import (
	"fmt"

	"github.com/auchon/auchon-api/lib/errors"
)

type InvalidLocationKindError struct{ errors.Base }

func NewInvalidLocationKindError(kind string) InvalidLocationKindError {
	return InvalidLocationKindError{errors.Base{
		Message: fmt.Sprintf("%q is not a kind of location (place, area, container)", kind),
		Context: map[string]any{"kind": kind},
	}}
}

func (InvalidLocationKindError) IsValidationError() bool { return true }
