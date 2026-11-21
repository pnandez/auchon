package create_location_test

import (
	"context"
	"testing"
	"time"

	"github.com/auchon/auchon-api/internal/inventory/application/create_location"
	"github.com/auchon/auchon-api/internal/inventory/domain"
	inventorytest "github.com/auchon/auchon-api/internal/inventory/infra/test"
	"github.com/auchon/auchon-api/lib/clock"
	"github.com/auchon/auchon-api/lib/test"
)

func TestCreateLocation(t *testing.T) {
	t.Run("Put a metal shelf in the basement", func(t *testing.T) {
		reset()
		givenTheBasement()
		whenCreated(shelfIn(basementID))
		thenSaved(t, domain.LocationDTO{ID: shelfID, Name: "Metal shelf", Kind: "area", ParentID: basementID})
		thenPublished(t, domain.LocationCreatedEventName)
	})

	t.Run("A shelf needs somewhere to be", func(t *testing.T) {
		reset()
		givenTheBasement()
		whenCreated(shelfIn(atticID))
		thenRejected[domain.ParentLocationNotFoundError](t)
		thenNothingSaved(t)
		thenPublished(t)
	})

	t.Run("A removed parent does not count", func(t *testing.T) {
		reset()
		givenTheBasement()
		givenTheAtticWasRemoved()
		whenCreated(shelfIn(atticID))
		thenRejected[domain.ParentLocationNotFoundError](t)
	})

	t.Run("A kind that does not exist is rejected", func(t *testing.T) {
		reset()
		givenTheBasement()
		cmd := shelfIn(basementID)
		cmd.Kind = "drawer"
		whenCreated(cmd)
		thenRejected[domain.InvalidLocationKindError](t)
		thenNothingSaved(t)
	})

	t.Run("A parent id that is not a UUID is rejected", func(t *testing.T) {
		reset()
		givenTheBasement()
		whenCreated(shelfIn("basement"))
		thenRejected[domain.InvalidLocationIDError](t)
	})

	t.Run("A location id that is not a UUID is rejected", func(t *testing.T) {
		reset()
		givenTheBasement()
		cmd := shelfIn(basementID)
		cmd.LocationID = "shelf"
		whenCreated(cmd)
		thenRejected[domain.InvalidLocationIDError](t)
	})
}

// ---- steps ----

const (
	basementID = "01937a4e-0000-7000-8000-000000000001"
	shelfID    = "01937a4e-0000-7000-8000-000000000002"
	atticID    = "01937a4e-0000-7000-8000-000000000009"
)

var (
	ctx     = context.Background()
	repo    = inventorytest.NewLocationRepositoryMock()
	events  = test.NewEventPublisherSpy()
	handler = create_location.NewHandler(domain.NewLocationFinder(repo), domain.NewLocationStorer(repo), events,
		clock.NewStub(time.Date(2026, 11, 21, 18, 30, 0, 0, time.UTC)))
	err error
)

func reset() { repo.Reset(); events.Reset(); err = nil }

func shelfIn(parentID string) create_location.CreateLocationCommand {
	return create_location.CreateLocationCommand{LocationID: shelfID, ParentID: parentID, Name: "Metal shelf", Kind: "area"}
}

func givenTheBasement() {
	repo.WithLocation(inventorytest.NewLocationBuilder().WithID(basementID).DTO())
}

func givenTheAtticWasRemoved() {
	repo.WithLocation(inventorytest.NewLocationBuilder().WithID(atticID).WithName("Attic").Removed().DTO())
}

func whenCreated(cmd create_location.CreateLocationCommand) { err = handler.Execute(ctx, cmd) }

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
