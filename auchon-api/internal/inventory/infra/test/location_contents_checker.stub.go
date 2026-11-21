package inventorytest

import "context"

// LocationContentsCheckerStub reports contents for the ids it was given.
type LocationContentsCheckerStub struct{ withContents map[string]bool }

func NewLocationContentsCheckerStub() *LocationContentsCheckerStub {
	return &LocationContentsCheckerStub{withContents: map[string]bool{}}
}

// Reset forgets every location with contents, for the next scenario.
func (s *LocationContentsCheckerStub) Reset() { s.withContents = map[string]bool{} }

func (s *LocationContentsCheckerStub) WithContents(locationID string) *LocationContentsCheckerStub {
	s.withContents[locationID] = true
	return s
}

func (s *LocationContentsCheckerStub) HasContents(_ context.Context, locationID string) (bool, error) {
	return s.withContents[locationID], nil
}
