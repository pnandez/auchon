package remove_location

import (
	"fmt"

	"github.com/auchon/auchon-api/lib/errors"
)

// LocationNotEmptyError: a location that still holds something cannot be removed.
type LocationNotEmptyError struct{ errors.Base }

func NewLocationNotEmptyError(id string) LocationNotEmptyError {
	return LocationNotEmptyError{errors.Base{
		Message: fmt.Sprintf("location %s still holds things; empty it first", id),
		Context: map[string]any{"locationId": id},
	}}
}

func (LocationNotEmptyError) IsValidationError() bool { return true }
