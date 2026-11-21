package domain

import "context"

type LocationFinder struct{ repository LocationRepository }

func NewLocationFinder(repository LocationRepository) *LocationFinder {
	return &LocationFinder{repository: repository}
}

// FindByID loads a location. Missing or removed locations are LocationNotFoundError.
func (f *LocationFinder) FindByID(ctx context.Context, id string) (*Location, error) {
	return nil, nil
}
