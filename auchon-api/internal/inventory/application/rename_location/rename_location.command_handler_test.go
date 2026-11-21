package rename_location_test

import (
	"context"
	"testing"
	"time"

	"github.com/auchon/auchon-api/internal/inventory/application/rename_location"
	"github.com/auchon/auchon-api/internal/inventory/domain"
	"github.com/auchon/auchon-api/internal/inventory/domain/name"
	inventorytest "github.com/auchon/auchon-api/internal/inventory/infra/test"
	"github.com/auchon/auchon-api/lib/clock"
	"github.com/auchon/auchon-api/lib/test"
)

func TestRenameLocation(t *testing.T) {
	t.Run("Rename a box", func(t *testing.T) {
		reset()
		givenBox1InTheBasement()
		whenRenamed(boxID, "Christmas lights")
		thenSaved(t, domain.LocationDTO{ID: boxID, Name: "Christmas lights", Kind: "container", ParentID: basementID})
		thenPublished(t, domain.LocationRenamedEventName)
	})

	t.Run("Renaming to the same name publishes nothing", func(t *testing.T) {
		reset()
		givenBox1InTheBasement()
		whenRenamed(boxID, "Box 1")
		thenSucceeded(t)
		thenPublished(t)
	})

	t.Run("A box that does not exist is not found", func(t *testing.T) {
		reset()
		givenBox1InTheBasement()
		whenRenamed(missingID, "X")
		thenRejected[domain.LocationNotFoundError](t)
	})

	t.Run("A box needs a name", func(t *testing.T) {
		reset()
		givenBox1InTheBasement()
		whenRenamed(boxID, " ")
		thenRejected[name.NameRequiredError](t)
		thenNothingSaved(t)
	})

	t.Run("An id that is not a UUID is rejected", func(t *testing.T) {
		reset()
		givenBox1InTheBasement()
		whenRenamed("box", "X")
		thenRejected[domain.InvalidLocationIDError](t)
	})
}

// ---- steps ----

const (
	basementID = "01937a4e-0000-7000-8000-000000000001"
	boxID      = "01937a4e-0000-7000-8000-000000000004"
	missingID  = "01937a4e-0000-7000-8000-000000000099"
)

var (
	ctx     = context.Background()
	repo    = inventorytest.NewLocationRepositoryMock()
	events  = test.NewEventPublisherSpy()
	handler = rename_location.NewHandler(domain.NewLocationFinder(repo), domain.NewLocationStorer(repo), events,
		clock.NewStub(time.Date(2026, 11, 21, 18, 30, 0, 0, time.UTC)))
	err error
)

func reset() { repo.Reset(); events.Reset(); err = nil }

func givenBox1InTheBasement() {
	repo.WithLocation(
		inventorytest.NewLocationBuilder().WithID(basementID).DTO(),
		inventorytest.NewLocationBuilder().WithID(boxID).WithName("Box 1").Inside(basementID).DTO())
}

func whenRenamed(id, n string) {
	err = handler.Execute(ctx, rename_location.RenameLocationCommand{LocationID: id, Name: n})
}

func thenSucceeded(t *testing.T) { t.Helper(); test.AssertNoError(t, err) }

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
