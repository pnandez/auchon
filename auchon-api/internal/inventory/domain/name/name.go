// Package name is the value object for the name of something in the inventory (a location or
// an item): trimmed, required, at most MaxLength characters.
package name

const MaxLength = 100

type Name struct{ value string }

// New trims s and validates it.
func New(s string) (Name, error) { return Name{}, nil }

// Hydrate rebuilds a stored name without validation. Never use it in tests.
func Hydrate(s string) Name { return Name{value: s} }

func (n Name) String() string         { return "" }
func (n Name) Equals(other Name) bool { return false }
