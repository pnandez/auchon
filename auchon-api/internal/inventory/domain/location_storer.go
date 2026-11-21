package domain

import (
	"context"
	"slices"
)

type LocationStorer struct{ repository LocationRepository }

func NewLocationStorer(repository LocationRepository) *LocationStorer {
	return &LocationStorer{repository: repository}
}

// Store persists the location after checking the collection invariants: a new location's id is
// not taken (LocationAlreadyExistsError), and no location ends up inside itself or one of the
// locations it contains (LocationCannotContainItselfError).
func (s *LocationStorer) Store(ctx context.Context, location *Location) error {
	if location.isNew {
		existing, err := s.repository.FindByID(ctx, location.id)
		if err != nil {
			return err
		}
		if existing != nil {
			return NewLocationAlreadyExistsError(location.id)
		}
	}
	if location.parentID != "" {
		ancestors, err := s.repository.AncestorIDs(ctx, location.parentID)
		if err != nil {
			return err
		}
		if location.parentID == location.id || slices.Contains(ancestors, location.id) {
			return NewLocationCannotContainItselfError(location.id)
		}
	}
	if err := s.repository.Save(ctx, location.dto()); err != nil {
		return err
	}
	location.isNew = false
	return nil
}
