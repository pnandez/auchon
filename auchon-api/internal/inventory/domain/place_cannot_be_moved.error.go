package domain

import (
	"fmt"

	"github.com/auchon/auchon-api/lib/errors"
)

type PlaceCannotBeMovedError struct{ errors.Base }

func NewPlaceCannotBeMovedError(id string) PlaceCannotBeMovedError {
	return PlaceCannotBeMovedError{errors.Base{Message: fmt.Sprintf("place %s cannot be moved inside another location", id), Context: map[string]any{"locationId": id}}}
}

func (PlaceCannotBeMovedError) IsValidationError() bool { return true }
