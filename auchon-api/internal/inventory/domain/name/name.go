// Package name is the value object for the name of something in the inventory (a location or
// an item): trimmed, required, at most MaxLength characters.
package name

import (
	"strings"
	"unicode/utf8"
)

const MaxLength = 100

type Name struct{ value string }

// New trims s and validates it.
func New(s string) (Name, error) {
	value := strings.TrimSpace(s)
	if value == "" {
		return Name{}, NewNameRequiredError()
	}
	if length := utf8.RuneCountInString(value); length > MaxLength {
		return Name{}, NewNameTooLongError(length)
	}
	return Name{value: value}, nil
}

// Hydrate rebuilds a stored name without validation. Never use it in tests.
func Hydrate(s string) Name { return Name{value: s} }

func (n Name) String() string         { return n.value }
func (n Name) Equals(other Name) bool { return n.value == other.value }
