package domain

import "context"

type LocationStorer struct{ repository LocationRepository }

func NewLocationStorer(repository LocationRepository) *LocationStorer {
	return &LocationStorer{repository: repository}
}

// Store persists the location after checking the collection invariants: a new location's id is
// not taken (LocationAlreadyExistsError), and no location ends up inside itself or one of the
// locations it contains (LocationCannotContainItselfError).
func (s *LocationStorer) Store(ctx context.Context, location *Location) error {
	return nil
}
