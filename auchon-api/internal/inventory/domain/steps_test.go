package domain_test

// Shared by location_test.go, location_finder_test.go and location_storer_test.go, which are
// one test package: the doubles, the result of the last step and the steps they all use.

import (
	"context"
	"testing"
	"time"

	"github.com/auchon/auchon-api/internal/inventory/domain"
	inventorytest "github.com/auchon/auchon-api/internal/inventory/infra/test"
	"github.com/auchon/auchon-api/lib/event"
	"github.com/auchon/auchon-api/lib/test"
)

const (
	basementID = "01937a4e-0000-7000-8000-000000000001"
	shelfAID   = "01937a4e-0000-7000-8000-000000000002"
	shelfBID   = "01937a4e-0000-7000-8000-000000000003"
	toolsID    = "01937a4e-0000-7000-8000-000000000004"
	screwsID   = "01937a4e-0000-7000-8000-000000000005"
	garageID   = "01937a4e-0000-7000-8000-000000000006"
)

var (
	ctx    = context.Background()
	at     = time.Date(2026, 11, 21, 18, 30, 0, 0, time.UTC)
	repo   = inventorytest.NewLocationRepositoryMock()
	finder = domain.NewLocationFinder(repo)
	storer = domain.NewLocationStorer(repo)

	location *domain.Location
	kind     domain.Kind
	err      error
)

func reset() { repo.Reset(); location = nil; kind = ""; err = nil }

// Locations of the basement: Basement > Shelf A > Tools > Screws, and Basement > Shelf B.
func basement() *inventorytest.LocationBuilder {
	return inventorytest.NewLocationBuilder().WithID(basementID)
}

func shelfA() *inventorytest.LocationBuilder {
	return inventorytest.NewLocationBuilder().WithID(shelfAID).WithName("Shelf A").WithKind(domain.KindArea).Inside(basementID)
}

func shelfB() *inventorytest.LocationBuilder {
	return inventorytest.NewLocationBuilder().WithID(shelfBID).WithName("Shelf B").WithKind(domain.KindArea).Inside(basementID)
}

func tools() *inventorytest.LocationBuilder {
	return inventorytest.NewLocationBuilder().WithID(toolsID).WithName("Tools").Inside(shelfAID)
}

func screws() *inventorytest.LocationBuilder {
	return inventorytest.NewLocationBuilder().WithID(screwsID).WithName("Screws").Inside(toolsID)
}

func givenStored(locations ...*inventorytest.LocationBuilder) {
	for _, l := range locations {
		repo.WithLocation(l.DTO())
	}
}

func thenSucceeded(t *testing.T) { t.Helper(); test.AssertNoError(t, err) }

func thenRejected[E error](t *testing.T) {
	t.Helper()
	test.AssertErrorIsOfType[E](t, err)
}

// thenItIs asserts the location from the last step has exactly these fields.
func thenItIs(t *testing.T, want domain.LocationDTO) {
	t.Helper()
	test.AssertNoError(t, err)
	test.AssertEqual(t, location.ID(), want.ID)
	test.AssertEqual(t, location.Name().String(), want.Name)
	test.AssertEqual(t, string(location.Kind()), want.Kind)
	test.AssertEqual(t, location.ParentID(), want.ParentID)
	test.AssertEqual(t, location.IsRemoved(), want.Removed)
}

func thenNothingRecorded(t *testing.T) {
	t.Helper()
	if events := location.PullEvents(); len(events) != 0 {
		t.Errorf("recorded %v, want no events", events)
	}
}

// recordedOnly asserts the location recorded exactly one event, of type E, about itself and
// at the scenario's time, and returns it.
func recordedOnly[E event.Event](t *testing.T) E {
	t.Helper()
	events := location.PullEvents()
	if len(events) != 1 {
		t.Fatalf("recorded %d events, want 1: %v", len(events), events)
	}
	e, ok := events[0].(E)
	if !ok {
		t.Fatalf("recorded %T", events[0])
	}
	if e.AggregateID() != location.ID() || !e.OccurredAt().Equal(at) {
		t.Errorf("event aggregate %q at %v, want %q at %v", e.AggregateID(), e.OccurredAt(), location.ID(), at)
	}
	return e
}
