package domain

import (
	"fmt"

	"github.com/auchon/auchon-api/lib/errors"
)

type LocationNotFoundError struct{ errors.Base }

func NewLocationNotFoundError(id string) LocationNotFoundError {
	return LocationNotFoundError{errors.Base{Message: fmt.Sprintf("location %s not found", id), Context: map[string]any{"locationId": id}}}
}

func (LocationNotFoundError) IsNotFound() bool { return true }
