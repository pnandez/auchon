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
func ParseKind(s string) (Kind, error) {
	switch k := Kind(s); k {
	case KindPlace, KindArea, KindContainer:
		return k, nil
	default:
		return "", NewInvalidLocationKindError(s)
	}
}

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
	n, err := name.New(rawName)
	if err != nil {
		return nil, err
	}
	return created(id, n, KindPlace, "", at), nil
}

// CreateLocationIn adds an area or a container inside the location parentID.
func CreateLocationIn(id, rawName string, kind Kind, parentID string, at time.Time) (*Location, error) {
	if kind != KindArea && kind != KindContainer {
		return nil, NewInvalidLocationKindError(string(kind))
	}
	if err := guardParent(id, parentID); err != nil {
		return nil, err
	}
	n, err := name.New(rawName)
	if err != nil {
		return nil, err
	}
	return created(id, n, kind, parentID, at), nil
}

func created(id string, n name.Name, kind Kind, parentID string, at time.Time) *Location {
	l := &Location{id: id, name: n, kind: kind, parentID: parentID, isNew: true}
	l.Record(LocationCreated{event.NewBase(id, at), LocationCreatedPayload{
		LocationID: id, Name: n.String(), Kind: string(kind), ParentID: parentID,
	}})
	return l
}

// Hydrate rebuilds a stored location without validation. Only Finders call it, never tests.
func Hydrate(id, storedName string, kind Kind, parentID string, removed bool) *Location {
	return &Location{id: id, name: name.Hydrate(storedName), kind: kind, parentID: parentID, removed: removed}
}

// Rename gives the location a new name. Renaming to the same name changes nothing.
func (l *Location) Rename(rawName string, at time.Time) error {
	n, err := name.New(rawName)
	if err != nil {
		return err
	}
	if n.Equals(l.name) {
		return nil
	}
	from := l.name
	l.name = n
	l.Record(LocationRenamed{event.NewBase(l.id, at), LocationRenamedPayload{
		LocationID: l.id, From: from.String(), To: n.String(),
	}})
	return nil
}

// MoveTo puts the location inside parentID. Places cannot move; moving to the current parent
// changes nothing.
func (l *Location) MoveTo(parentID string, at time.Time) error {
	if l.kind == KindPlace {
		return NewPlaceCannotBeMovedError(l.id)
	}
	if err := guardParent(l.id, parentID); err != nil {
		return err
	}
	if parentID == l.parentID {
		return nil
	}
	from := l.parentID
	l.parentID = parentID
	l.Record(LocationMoved{event.NewBase(l.id, at), LocationMovedPayload{
		LocationID: l.id, From: from, To: parentID,
	}})
	return nil
}

// Remove takes the location out of the inventory.
func (l *Location) Remove(at time.Time) {
	l.removed = true
	l.Record(LocationRemoved{event.NewBase(l.id, at), LocationRemovedPayload{LocationID: l.id}})
}

func (l *Location) ID() string       { return l.id }
func (l *Location) Name() name.Name  { return l.name }
func (l *Location) Kind() Kind       { return l.kind }
func (l *Location) ParentID() string { return l.parentID }
func (l *Location) IsRemoved() bool  { return l.removed }

func (l *Location) dto() LocationDTO {
	return LocationDTO{ID: l.id, Name: l.name.String(), Kind: string(l.kind), ParentID: l.parentID, Removed: l.removed}
}

// guardParent checks what the location alone can know about its parent; the storer checks
// the deeper nesting.
func guardParent(id, parentID string) error {
	if parentID == "" {
		return NewLocationParentRequiredError(id)
	}
	if parentID == id {
		return NewLocationCannotContainItselfError(id)
	}
	return nil
}
