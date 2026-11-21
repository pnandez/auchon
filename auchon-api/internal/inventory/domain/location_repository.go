package domain

import "context"

// LocationDTO is how a location crosses the repository port.
type LocationDTO struct {
	ID       string
	Name     string
	Kind     string
	ParentID string // empty for places
	Removed  bool
}

type LocationRepository interface {
	Save(ctx context.Context, location LocationDTO) error
	// FindByID returns nil when the location does not exist (removed ones are returned).
	FindByID(ctx context.Context, id string) (*LocationDTO, error)
	// AncestorIDs returns the ids of the locations that contain id, nearest first.
	AncestorIDs(ctx context.Context, id string) ([]string, error)
}
