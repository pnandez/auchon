package domain

import "context"

// LocationContentsChecker tells whether a location still holds anything (other locations, and
// later items).
type LocationContentsChecker interface {
	HasContents(ctx context.Context, locationID string) (bool, error)
}
