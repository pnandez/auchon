package domain

import (
	"time"

	"github.com/auchon/auchon-api/internal/inventory/domain/name"
	"github.com/auchon/auchon-api/lib/event"
)

// Kind says what a location is. A place (Basement, Garage…) is a root; areas (a shelf, a
// corner) and containers (a box, a drawer) always sit inside another location.
type Kind string

const (
	KindPlace     Kind = "place"
	KindArea      Kind = "area"
	KindContainer Kind = "container"
)

// ParseKind returns the kind named s, or InvalidLocationKindError.
func ParseKind(s string) (Kind, error) { return "", nil }

// Location is the aggregate root for where things are kept.
type Location struct {
	event.Recorder

	id       string
	name     name.Name
	kind     Kind
	parentID string
	removed  bool
	isNew    bool
}

// CreatePlace adds a new top-level place, such as the basement.
func CreatePlace(id, rawName string, at time.Time) (*Location, error) {
	return nil, nil
}

// CreateLocationIn adds an area or a container inside the location parentID.
func CreateLocationIn(id, rawName string, kind Kind, parentID string, at time.Time) (*Location, error) {
	return nil, nil
}

// Hydrate rebuilds a stored location without validation. Only Finders call it, never tests.
func Hydrate(id, storedName string, kind Kind, parentID string, removed bool) *Location {
	return &Location{}
}

// Rename gives the location a new name. Renaming to the same name changes nothing.
func (l *Location) Rename(rawName string, at time.Time) error { return nil }

// MoveTo puts the location inside parentID. Places cannot move; moving to the current parent
// changes nothing.
func (l *Location) MoveTo(parentID string, at time.Time) error { return nil }

// Remove takes the location out of the inventory.
func (l *Location) Remove(at time.Time) {}

func (l *Location) ID() string       { return "" }
func (l *Location) Name() name.Name  { return name.Name{} }
func (l *Location) Kind() Kind       { return "" }
func (l *Location) ParentID() string { return "" }
func (l *Location) IsRemoved() bool  { return false }
