package move_location_test

import (
	"context"
	"testing"
	"time"

	"github.com/auchon/auchon-api/internal/inventory/application/move_location"
	"github.com/auchon/auchon-api/internal/inventory/domain"
	inventorytest "github.com/auchon/auchon-api/internal/inventory/infra/test"
	"github.com/auchon/auchon-api/lib/clock"
	"github.com/auchon/auchon-api/lib/test"
)

func TestMoveLocation(t *testing.T) {
	t.Run("Move a box to another shelf", func(t *testing.T) {
		reset()
		givenTheBasementWithTwoShelves()
		givenToolsOnShelfA()
		whenMoved(toolsID, shelfBID)
		thenSaved(t, domain.LocationDTO{ID: toolsID, Name: "Tools", Kind: "container", ParentID: shelfBID})
		thenPublished(t, domain.LocationMovedEventName)
	})

	t.Run("A box cannot go inside a box it contains", func(t *testing.T) {
		reset()
		givenTheBasementWithTwoShelves()
		givenToolsOnShelfA()
		givenScrewsInTools()
		whenMoved(toolsID, screwsID)
		thenRejected[domain.LocationCannotContainItselfError](t)
		thenNothingSaved(t)
		thenPublished(t)
	})

	t.Run("A place stays at the top", func(t *testing.T) {
		reset()
		givenTheBasementWithTwoShelves()
		givenTheGarage()
		whenMoved(garageID, basementID)
		thenRejected[domain.PlaceCannotBeMovedError](t)
		thenNothingSaved(t)
	})

	t.Run("A location that does not exist is not found", func(t *testing.T) {
		reset()
		givenTheBasementWithTwoShelves()
		whenMoved(missingID, shelfBID)
		thenRejected[domain.LocationNotFoundError](t)
	})

	t.Run("A box needs somewhere to go", func(t *testing.T) {
		reset()
		givenTheBasementWithTwoShelves()
		givenToolsOnShelfA()
		whenMoved(toolsID, missingID)
		thenRejected[domain.ParentLocationNotFoundError](t)
		thenNothingSaved(t)
	})

	t.Run("A location id that is not a UUID is rejected", func(t *testing.T) {
		reset()
		givenTheBasementWithTwoShelves()
		whenMoved("tools", shelfBID)
		thenRejected[domain.InvalidLocationIDError](t)
	})

	t.Run("A parent id that is not a UUID is rejected", func(t *testing.T) {
		reset()
		givenTheBasementWithTwoShelves()
		givenToolsOnShelfA()
		whenMoved(toolsID, "shelf-b")
		thenRejected[domain.InvalidLocationIDError](t)
	})
}

// ---- steps ----

const (
	basementID = "01937a4e-0000-7000-8000-000000000001"
	shelfAID   = "01937a4e-0000-7000-8000-000000000002"
	shelfBID   = "01937a4e-0000-7000-8000-000000000003"
	toolsID    = "01937a4e-0000-7000-8000-000000000004"
	screwsID   = "01937a4e-0000-7000-8000-000000000005"
	garageID   = "01937a4e-0000-7000-8000-000000000006"
	missingID  = "01937a4e-0000-7000-8000-000000000099"
)

var (
	ctx     = context.Background()
	repo    = inventorytest.NewLocationRepositoryMock()
	events  = test.NewEventPublisherSpy()
	handler = move_location.NewHandler(domain.NewLocationFinder(repo), domain.NewLocationStorer(repo), events,
		clock.NewStub(time.Date(2026, 11, 21, 18, 30, 0, 0, time.UTC)))
	err error
)

func reset() { repo.Reset(); events.Reset(); err = nil }

func givenTheBasementWithTwoShelves() {
	repo.WithLocation(
		inventorytest.NewLocationBuilder().WithID(basementID).DTO(),
		inventorytest.NewLocationBuilder().WithID(shelfAID).WithName("Shelf A").WithKind(domain.KindArea).Inside(basementID).DTO(),
		inventorytest.NewLocationBuilder().WithID(shelfBID).WithName("Shelf B").WithKind(domain.KindArea).Inside(basementID).DTO())
}

func givenToolsOnShelfA() {
	repo.WithLocation(inventorytest.NewLocationBuilder().WithID(toolsID).WithName("Tools").Inside(shelfAID).DTO())
}

func givenScrewsInTools() {
	repo.WithLocation(inventorytest.NewLocationBuilder().WithID(screwsID).WithName("Screws").Inside(toolsID).DTO())
}

func givenTheGarage() {
	repo.WithLocation(inventorytest.NewLocationBuilder().WithID(garageID).WithName("Garage").DTO())
}

func whenMoved(id, to string) {
	err = handler.Execute(ctx, move_location.MoveLocationCommand{LocationID: id, ParentID: to})
}

func thenSaved(t *testing.T, want domain.LocationDTO) {
	t.Helper()
	test.AssertNoError(t, err)
	repo.AssertSaved(t, want)
}

func thenRejected[E error](t *testing.T) {
	t.Helper()
	test.AssertErrorIsOfType[E](t, err)
}

func thenNothingSaved(t *testing.T) { t.Helper(); repo.AssertNothingSaved(t) }

func thenPublished(t *testing.T, names ...string) { t.Helper(); events.AssertPublished(t, names...) }
