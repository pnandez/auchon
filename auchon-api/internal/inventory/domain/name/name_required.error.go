package name

import "github.com/auchon/auchon-api/lib/errors"

type NameRequiredError struct{ errors.Base }

func NewNameRequiredError() NameRequiredError {
	return NameRequiredError{errors.Base{Message: "a name is required"}}
}

func (NameRequiredError) IsValidationError() bool { return true }
