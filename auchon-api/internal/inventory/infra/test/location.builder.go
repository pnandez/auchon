package inventorytest

import (
	"time"

	"github.com/auchon/auchon-api/internal/inventory/domain"
	"github.com/google/uuid"
)

var builtAt = time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC)

// LocationBuilder builds locations through the domain factories (never Hydrate). Defaults:
// a place called "Basement" with a fresh id.
type LocationBuilder struct {
	id       string
	name     string
	kind     domain.Kind
	parentID string
	removed  bool
}

func NewLocationBuilder() *LocationBuilder {
	return &LocationBuilder{id: uuid.Must(uuid.NewV7()).String(), name: "Basement", kind: domain.KindPlace}
}

func (b *LocationBuilder) WithID(id string) *LocationBuilder       { b.id = id; return b }
func (b *LocationBuilder) WithName(n string) *LocationBuilder      { b.name = n; return b }
func (b *LocationBuilder) Removed() *LocationBuilder               { b.removed = true; return b }
func (b *LocationBuilder) WithKind(k domain.Kind) *LocationBuilder { b.kind = k; return b }

// Inside makes it a container inside parentID (unless WithKind sets area).
func (b *LocationBuilder) Inside(parentID string) *LocationBuilder {
	b.parentID = parentID
	if b.kind == domain.KindPlace {
		b.kind = domain.KindContainer
	}
	return b
}

// Build returns the entity with its recorded events cleared.
func (b *LocationBuilder) Build() *domain.Location {
	var (
		l   *domain.Location
		err error
	)
	if b.kind == domain.KindPlace {
		l, err = domain.CreatePlace(b.id, b.name, builtAt)
	} else {
		l, err = domain.CreateLocationIn(b.id, b.name, b.kind, b.parentID, builtAt)
	}
	if err != nil {
		panic("LocationBuilder: " + err.Error())
	}
	if b.removed {
		l.Remove(builtAt)
	}
	l.PullEvents()
	return l
}

// DTO returns the stored form, for seeding a repository.
func (b *LocationBuilder) DTO() domain.LocationDTO {
	return domain.LocationDTO{ID: b.id, Name: b.name, Kind: string(b.kind), ParentID: b.parentID, Removed: b.removed}
}
