package domain

import "context"

type LocationFinder struct{ repository LocationRepository }

func NewLocationFinder(repository LocationRepository) *LocationFinder {
	return &LocationFinder{repository: repository}
}

// FindByID loads a location. Missing or removed locations are LocationNotFoundError.
func (f *LocationFinder) FindByID(ctx context.Context, id string) (*Location, error) {
	dto, err := f.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if dto == nil || dto.Removed {
		return nil, NewLocationNotFoundError(id)
	}
	return Hydrate(dto.ID, dto.Name, Kind(dto.Kind), dto.ParentID, dto.Removed), nil
}
