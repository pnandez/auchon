package name

import (
	"fmt"

	"github.com/auchon/auchon-api/lib/errors"
)

type NameTooLongError struct{ errors.Base }

func NewNameTooLongError(length int) NameTooLongError {
	return NameTooLongError{errors.Base{
		Message: fmt.Sprintf("a name can have at most %d characters, got %d", MaxLength, length),
		Context: map[string]any{"length": length, "max": MaxLength},
	}}
}

func (NameTooLongError) IsValidationError() bool { return true }
