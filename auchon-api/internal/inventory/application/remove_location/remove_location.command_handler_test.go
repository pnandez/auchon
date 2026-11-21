package remove_location_test

import (
	"context"
	"testing"
	"time"

	"github.com/auchon/auchon-api/internal/inventory/application/remove_location"
	"github.com/auchon/auchon-api/internal/inventory/domain"
	inventorytest "github.com/auchon/auchon-api/internal/inventory/infra/test"
	"github.com/auchon/auchon-api/lib/clock"
	"github.com/auchon/auchon-api/lib/test"
)

func TestRemoveLocation(t *testing.T) {
	t.Run("Remove an empty box", func(t *testing.T) {
		reset()
		givenAnOldBoxInTheBasement()
		whenRemoved(boxID)
		thenSaved(t, domain.LocationDTO{ID: boxID, Name: "Old box", Kind: "container", ParentID: basementID, Removed: true})
		thenPublished(t, domain.LocationRemovedEventName)
	})

	t.Run("A box that still holds things cannot be removed", func(t *testing.T) {
		reset()
		givenAnOldBoxInTheBasement()
		givenItHoldsThings(boxID)
		whenRemoved(boxID)
		thenRejected[remove_location.LocationNotEmptyError](t)
		thenNothingSaved(t)
		thenPublished(t)
	})

	t.Run("A box that does not exist is not found", func(t *testing.T) {
		reset()
		givenAnOldBoxInTheBasement()
		whenRemoved(missingID)
		thenRejected[domain.LocationNotFoundError](t)
	})

	t.Run("An id that is not a UUID is rejected", func(t *testing.T) {
		reset()
		givenAnOldBoxInTheBasement()
		whenRemoved("box")
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
	ctx      = context.Background()
	repo     = inventorytest.NewLocationRepositoryMock()
	contents = inventorytest.NewLocationContentsCheckerStub()
	events   = test.NewEventPublisherSpy()
	handler  = remove_location.NewHandler(domain.NewLocationFinder(repo), domain.NewLocationStorer(repo), contents, events,
		clock.NewStub(time.Date(2026, 11, 21, 18, 30, 0, 0, time.UTC)))
	err error
)

func reset() { repo.Reset(); contents.Reset(); events.Reset(); err = nil }

func givenAnOldBoxInTheBasement() {
	repo.WithLocation(
		inventorytest.NewLocationBuilder().WithID(basementID).DTO(),
		inventorytest.NewLocationBuilder().WithID(boxID).WithName("Old box").Inside(basementID).DTO())
}

func givenItHoldsThings(id string) { contents.WithContents(id) }

func whenRemoved(id string) {
	err = handler.Execute(ctx, remove_location.RemoveLocationCommand{LocationID: id})
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
